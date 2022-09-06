package auth

import (
	"context"
	"errors"
	"github.com/wolframdeus/exchange-rates-backend/configs"
	"github.com/wolframdeus/exchange-rates-backend/internal/db/models"
	"github.com/wolframdeus/exchange-rates-backend/internal/jwt"
	"github.com/wolframdeus/exchange-rates-backend/internal/repositories/refsessions"
	"github.com/wolframdeus/exchange-rates-backend/internal/services/users"
	"github.com/wolframdeus/exchange-rates-backend/internal/tg"
	"sort"
	"time"
)

type Jwt struct {
	Token     string
	ExpiresAt time.Time
}

type Result struct {
	// Токен доступа.
	AccessToken *Jwt
	// Токен для обновления токена доступа.
	RefreshToken *Jwt
}

type Service struct {
	// Сервис для работы с пользователями.
	uSrv *users.Service
	// Репозиторий для работы с сессиями.
	refRep *refsessions.Repository
}

// AuthenticateTg аутентифицирует пользователя по его параметрам запуска.
func (s *Service) AuthenticateTg(
	ctx context.Context,
	initData string,
	fp string,
) (*Result, error) {
	// Валидируем параметры запуска.
	if ok, err := tg.ValidateInitData(initData, configs.Tg.SecretKey); err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("init data is invalid")
	}

	// Парсим параметры запуска.
	p, err := tg.ParseInitData(initData)
	if err != nil {
		return nil, err
	}
	if p.User == nil {
		return nil, errors.New("user field is missing in init data")
	}

	// Получаем информацию о пользователе по его идентификатору Telegram.
	u, err := s.uSrv.FindByTelegramUid(ctx, p.User.Id)
	if err != nil {
		return nil, err
	}

	// Пользователь не существовал, поэтому его необходимо создать.
	if u == nil {
		// В противном случае просто создадим пользователя.
		u, err = s.uSrv.CreateByTgUid(ctx, p.User.Id, p.User.LanguageCode)
		if err != nil {
			return nil, err
		}

		// Создаем токены для пользователя.
		ut, err := jwt.CreateUserToken(u.Id, u.Lang)
		if err != nil {
			return nil, err
		}

		// Создаем сессию пользователя.
		_, err = s.refRep.Create(
			ctx,
			u.Id,
			ut.RefreshToken.Token,
			ut.AccessToken.Token,
			fp,
			ut.RefreshToken.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}

		return resultFromUserToken(ut), nil
	}

	// Если пользователь существовал, то нужно совершать дополнительные
	// операции связанные с session-менджментом.

	// Создаем токены для пользователя.
	ut, err := jwt.CreateUserToken(u.Id, u.Lang)
	if err != nil {
		return nil, err
	}

	// Получаем список всех сессий.
	sess, err := s.refRep.FindByUserId(ctx, u.Id)
	if err != nil {
		return nil, err
	}

	// Пытаемся найти предыдущую сессию с таким же fingerprint-ом.
	// Если такая сессия существует, то мы её просто инвалидируем и
	// обновим в БД.
	var prevSession *models.RefreshSession

	for _, session := range sess {
		if session.Fingerprint == fp {
			prevSession = session
			break
		}
	}

	// Предыдущая сессия была найдена. Инвалидируем её.
	if prevSession != nil {
		// TODO: Инвалидировать access token, который был до этого в сессии.
		// Обновляем информацию о сессии.
		if err := s.refRep.RefreshById(
			ctx,
			prevSession.Id,
			ut.RefreshToken.Token,
			ut.RefreshToken.ExpiresAt,
			ut.AccessToken.Token,
		); err != nil {
			return nil, err
		}

		// Переназначаем дату создания предыдущей сессии для возможных
		// дальнейших алгоритмов, которые эти сессии сортируем по дате
		// создания.
		prevSession.CreatedAt = time.Now()
	} else {
		newSess, err := s.refRep.Create(
			ctx,
			u.Id,
			ut.RefreshToken.Token,
			ut.AccessToken.Token,
			fp,
			ut.RefreshToken.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}

		// Добавляем новосозданную сессию в список известных.
		sess = append(sess, newSess)
	}

	// Мы не допускаем одновременного наличия более 5 сессий.
	if len(sess) > 5 {
		// Сортируем сессии по дате убывания их создания.
		sort.SliceStable(sess, func(i, j int) bool {
			return sess[i].CreatedAt.After(sess[j].CreatedAt)
		})

		// Получаем список идентификаторов тех сессий, которые нужно
		// инвалидировать.
		dropSessIds := make([]models.RefreshSessionId, len(sess)-5)
		for i, session := range sess[5:] {
			dropSessIds[i] = session.Id
		}

		// TODO: Инвалидировать access token-ы в сессиях.
		if _, err := s.refRep.DeleteByIds(ctx, dropSessIds); err != nil {
			return nil, err
		}
	}

	return resultFromUserToken(ut), nil
}

// RefreshSession обновляет сессию по указанному токену.
func (s *Service) RefreshSession(
	ctx context.Context,
	token string,
	fingerprint string,
) (*Result, error) {
	// Проверяем, валиден ли сам токен.
	if err := jwt.DecodeRefreshToken(token); err != nil {
		return nil, err
	}

	// Находим сессию.
	sess, err := s.refRep.FindByRefreshToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if sess == nil {
		return nil, errors.New("session not found")
	}

	// В случае, если обновление вызывается иным клиентом, необходимо сбросить
	// все сессии.
	if sess.Fingerprint != fingerprint {
		// TODO: Сбросить все сессии.

		// Удаляем сессии из БД, чтобы ими невозможно было более воспользоваться.
		_, err := s.refRep.DeleteByUserId(ctx, sess.UserId)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("unauthorized access detected")
	}

	// Находим пользователя-владельца сессии.
	u, err := s.uSrv.FindById(ctx, sess.UserId)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("user not found")
	}

	// Создаем токены для пользователя.
	ut, err := jwt.CreateUserToken(u.Id, u.Lang)
	if err != nil {
		return nil, err
	}

	// Обновляем информацию о сессии.
	if err := s.refRep.RefreshById(
		ctx,
		sess.Id,
		ut.RefreshToken.Token,
		ut.RefreshToken.ExpiresAt,
		ut.AccessToken.Token,
	); err != nil {
		return nil, err
	}
	return resultFromUserToken(ut), nil
}

// NewService создает указатель на новый экземпляр Service.
func NewService(uSrv *users.Service, refRep *refsessions.Repository) *Service {
	return &Service{uSrv: uSrv, refRep: refRep}
}

// Создает указатель на Result из результата авторизации пользователя.
func resultFromUserToken(ut *jwt.UserToken) *Result {
	return &Result{
		AccessToken: &Jwt{
			Token:     ut.AccessToken.Token,
			ExpiresAt: ut.AccessToken.ExpiresAt,
		},
		RefreshToken: &Jwt{
			Token:     ut.RefreshToken.Token,
			ExpiresAt: ut.RefreshToken.ExpiresAt,
		},
	}
}
