package services_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/core/locale"
	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
)

// writeTranslations writes a translation file for the locale and gives its path.
// The language of a file comes from its name.
func writeTranslations(tb testing.TB, localeCode, greeting string) string {
	tb.Helper()

	path := filepath.Join(tb.TempDir(), localeCode+".json")
	content := `[{"id": "greeting", "translation": "` + greeting + `"}]`

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		tb.Fatalf("write %s: %v", path, err)
	}

	return path
}

func TestLocale_TranslateLabelAndDateLocation(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	english := writeTranslations(t, "en-GB", "Hello")
	german := writeTranslations(t, "de-DE", "Hallo")

	app := harness.Boot(t, []dingo.Module{new(locale.Module)}, harness.Config{
		Files: map[string]string{
			"config.yml": "flamingo.debug.mode: false\n" +
				"core.locale.locale: en-GB\n" +
				"core.locale.translationFiles: [\"" + english + "\", \"" + german + "\"]\n" +
				"core.locale.date.location: Europe/Berlin\n",
		},
	})

	injector := app.Injector(t)

	t.Run("translate", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		instance, err := injector.GetInstance(new(domain.TranslationService))
		require.NoError(t, err)

		service, ok := instance.(domain.TranslationService)
		require.True(t, ok)

		assert.Equal(t, "Hello", service.Translate("greeting", "", "en-GB", 1, nil))
		assert.Equal(t, "Hallo", service.Translate("greeting", "", "de-DE", 1, nil))
	})

	t.Run("label", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		instance, err := injector.GetInstance(new(application.LabelService))
		require.NoError(t, err)

		service, ok := instance.(*application.LabelService)
		require.True(t, ok)

		assert.Equal(t, "Hello", service.NewLabel("greeting").String())
	})

	t.Run("date location", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		instance, err := injector.GetInstance(new(application.DateTimeServiceInterface))
		require.NoError(t, err)

		service, ok := instance.(application.DateTimeServiceInterface)
		require.True(t, ok)

		// 22:30 on the 1st of June in UTC is 00:30 on the 2nd in Berlin, and 23:30 on the 1st in London
		formatter, err := service.GetDateTimeFormatter(time.Date(2026, time.June, 1, 22, 30, 0, 0, time.UTC))
		require.NoError(t, err)

		assert.Equal(t, "02 Jun 2026", formatter.FormatToLocalDate())
		assert.Equal(t, "00:30", formatter.FormatLocale("15:04"))
	})
}
