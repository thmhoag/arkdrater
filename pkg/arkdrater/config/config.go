package config

type Config struct {
	Server struct {
		Port string `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	DynamicConfig struct {
		Multipliers struct {
			TamingSpeedMultiplier               float32 `yaml:"TamingSpeedMultiplier" envconfig:"arkdrater_TamingSpeedMultiplier"`
			HarvestAmountMultiplier             float32 `yaml:"HarvestAmountMultiplier" envconfig:"arkdrater_HarvestAmountMultiplier"`
			XPMultiplier                        float32 `yaml:"XPMultiplier" envconfig:"arkdrater_XPMultiplier"`
			MatingIntervalMultiplier            float32 `yaml:"MatingIntervalMultiplier" envconfig:"arkdrater_MatingIntervalMultiplier"`
			BabyMatureSpeedMultiplier           float32 `yaml:"BabyMatureSpeedMultiplier" envconfig:"arkdrater_BabyMatureSpeedMultiplier"`
			EggHatchSpeedMultiplier             float32 `yaml:"EggHatchSpeedMultiplier" envconfig:"arkdrater_EggHatchSpeedMultiplier"`
			BabyFoodConsumptionSpeedMultiplier  float32 `yaml:"BabyFoodConsumptionSpeedMultiplier" envconfig:"arkdrater_BabyFoodConsumptionSpeedMultiplier"`
			CropGrowthSpeedMultiplier           float32 `yaml:"CropGrowthSpeedMultiplier" envconfig:"arkdrater_CropGrowthSpeedMultiplier"`
			MatingSpeedMultiplier               float32 `yaml:"MatingSpeedMultiplier" envconfig:"arkdrater_MatingSpeedMultiplier"`
			BabyCuddleIntervalMultiplier        float32 `yaml:"BabyCuddleIntervalMultiplier" envconfig:"arkdrater_BabyCuddleIntervalMultiplier"`
			BabyImprintAmountMultiplier         float32 `yaml:"BabyImprintAmountMultiplier" envconfig:"arkdrater_BabyImprintAmountMultiplier"`
			CustomRecipeEffectivenessMultiplier float32 `yaml:"BabyImprintAmountMultiplier" envconfig:"arkdrater_BabyImprintAmountMultiplier"`
			HexagonRewardMultiplier             float32 `yaml:"HexagonRewardMultiplier" envconfig:"arkdrater_HexagonRewardMultiplier"`
		}
		TributeItemExpirationSeconds int64 `yaml:"TributeItemExpirationSeconds" envconfig:"arkdrater_TributeItemExpirationSeconds"`
		TributeDinoExpirationSeconds int64 `yaml:"TributeDinoExpirationSeconds" envconfig:"arkdrater_TributeDinoExpirationSeconds"`
		EnableFullDump               bool  `yaml:"EnableFullDump" envconfig:"arkdrater_EnableFullDump"`
		GUseServerNetSpeedCheck      bool  `yaml:"GUseServerNetSpeedCheck" envconfig:"arkdrater_GUseServerNetSpeedCheck"`
		bUseAlarmNotifications       bool  `yaml:"bUseAlarmNotifications" envconfig:"arkdrater_BbUseAlarmNotifications"`
	}
}
