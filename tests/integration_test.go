package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	// blank import for app wiring registration
	_ "cosmossdk.io/x/accounts"
	_ "cosmossdk.io/x/auth"
	_ "cosmossdk.io/x/auth/tx/config"
	_ "cosmossdk.io/x/bank"
	_ "cosmossdk.io/x/consensus"
	_ "cosmossdk.io/x/mint"
	_ "cosmossdk.io/x/staking"
	_ "github.com/cosmos/cosmos-sdk/x/genutil"
	_ "github.com/cosmosregistry/example/module"

	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/testutil/configurator"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	"github.com/cosmosregistry/example"
	"github.com/cosmosregistry/example/keeper"
)

// ExampleModule is a configurator.ModuleOption that add the example module to the app config.
var ExampleModule = func() configurator.ModuleOption {
	return func(config *configurator.Config) {
		config.ModuleConfigs[example.ModuleName] = &appv1alpha1.ModuleConfig{
			Name:   example.ModuleName,
			Config: appconfig.WrapAny(&example.Module{}),
		}
	}
}

func TestIntegration(t *testing.T) {
	t.Parallel()

	logger := log.NewTestLogger(t)
	appConfig := depinject.Configs(
		configurator.NewAppConfig(
			configurator.AuthModule(),
			configurator.BankModule(),
			configurator.StakingModule(),
			configurator.TxModule(),
			configurator.ConsensusModule(),
			configurator.GenutilModule(),
			configurator.MintModule(),
			configurator.AccountsModule(),
			ExampleModule(),
			configurator.WithCustomInitGenesisOrder(
				"accounts",
				"auth",
				"bank",
				"staking",
				"mint",
				"genutil",
				"consensus",
				example.ModuleName,
			),
		),
		depinject.Supply(logger))

	var keeper keeper.Keeper
	app, err := simtestutil.Setup(appConfig, &keeper)
	require.NoError(t, err)
	require.NotNil(t, app) // use the app or the keeper for running integration tests
}
