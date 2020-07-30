package command

import (
	"errors"
	"flag"
	"io/ioutil"
	"testing"

	"github.com/jonyhy96/kt-connect/fake/kt"

	"github.com/jonyhy96/kt-connect/fake/kt/action"
	"github.com/jonyhy96/kt-connect/pkg/kt/options"
	"github.com/golang/mock/gomock"
	"github.com/urfave/cli"
)

func Test_meshCommand(t *testing.T) {

	ctl := gomock.NewController(t)
	fakeKtCli := kt.NewMockCliInterface(ctl)
	mockAction := action.NewMockActionInterface(ctl)

	mockAction.EXPECT().Mesh(gomock.Eq("service"), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	cases := []struct {
		testArgs               []string
		skipFlagParsing        bool
		useShortOptionHandling bool
		expectedErr            error
	}{
		{testArgs: []string{"mesh", "service", "--expose", "8080"}, skipFlagParsing: false, useShortOptionHandling: false, expectedErr: nil},
		{testArgs: []string{"mesh", "service"}, skipFlagParsing: false, useShortOptionHandling: false, expectedErr: errors.New("-expose is required")},
		{testArgs: []string{"mesh"}, skipFlagParsing: false, useShortOptionHandling: false, expectedErr: errors.New("mesh target is required")},
	}

	for _, c := range cases {

		app := &cli.App{Writer: ioutil.Discard}
		set := flag.NewFlagSet("test", 0)
		_ = set.Parse(c.testArgs)

		context := cli.NewContext(app, set, nil)

		opts := options.NewDaemonOptions()
		opts.Debug = true
		command := newMeshCommand(fakeKtCli, opts, mockAction)
		err := command.Run(context)

		if c.expectedErr != nil {
			if err.Error() != c.expectedErr.Error() {
				t.Errorf("expected %v but is %v", c.expectedErr, err)
			}
		} else if err != c.expectedErr {
			t.Errorf("expected %v but is %v", c.expectedErr, err)
		}

	}

}
