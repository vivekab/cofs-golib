package goliblocalise

import (
	"context"
	"time"

	goliblogging "github.com/vivekab/golib/pkg/logging"
	golibtypes "github.com/vivekab/golib/pkg/types"
	grpcRoot "github.com/vivekab/golib/protobuf/protoroot"
)

type localizationStore struct {
	store          map[string]LocaliseAttributes // language:localisation_key->value ex : en:EC_WEBHOOK_ERROR->value
	refreshTimer   *time.Ticker
	serviceName    golibtypes.ServiceName
	localiseClient LocaliseService
}

type LocaliseAttributes struct {
	Value    string
	HttpCode int32
}

func newStore(refreshInterval int, serviceName golibtypes.ServiceName, client LocaliseService) *localizationStore {
	return &localizationStore{
		refreshTimer:   time.NewTicker(time.Duration(refreshInterval) * time.Second),
		store:          map[string]LocaliseAttributes{},
		localiseClient: client,
		serviceName:    serviceName,
	}
}

func getLocaliseStoreMapKey(lang grpcRoot.Language, key string) string {
	return string(golibtypes.LanguageFromProto[lang]) + ":" + key
}

func (l *localizationStore) Insert(language grpcRoot.Language, key string, value LocaliseAttributes) {
	mapkey := getLocaliseStoreMapKey(language, key)
	l.store[mapkey] = value
}

func (l *localizationStore) Get(key string, language grpcRoot.Language) (string, bool) {
	mapkey := getLocaliseStoreMapKey(language, key)
	value, ok := l.store[mapkey]
	if ok {
		return value.Value, true
	}
	return "", false
}

// func (l *localizationStore) loop() {
// 	for range l.refreshTimer.C {
// 		l.updateStore()
// 	}
// }

// func (l *localizationStore) updateStore() {
// 	l.loadTranslations()
// }

func (l *localizationStore) Start() {
	err := l.loadTranslations()
	if err != nil {
		goliblogging.Warn("failed to get translations : %s", err.Error())
	}
	// go l.loop()
}

func (l *localizationStore) Stop() {}

func (l *localizationStore) loadTranslations() error {
	if l.localiseClient == nil {
		return nil
	}

	for _, file := range []string{"global", string(l.serviceName)} {
		resp, err := l.localiseClient.ListLocaliseTranslation(context.Background(), &ListLocaliseTranslationRequest{FileName: file})
		if err != nil {
			return err
		}

		for _, translation := range resp {
			l.Insert(golibtypes.LanguageToProto[golibtypes.Language(translation.Language)], translation.Key, LocaliseAttributes{
				Value:    translation.Value,
				HttpCode: int32(translation.HTTPStatusCode),
			})
		}
	}

	goliblogging.Info("translations update successful")
	return nil
}
