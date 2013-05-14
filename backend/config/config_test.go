package config

import "testing"

func setup() []*Config {
	cfg := make([]*Config, 4)
	cfg[0] = NewConfig("testcfg0.cfg")
	cfg[1] = NewConfig("testcfg1.cfg")
	cfg[2] = NewConfig("testcfg2.cfg")
	cfg[3] = NewConfig("testcfg3.cfg")

	cfg[1].AddCfgValue("NICK", "Pelle")
	cfg[1].AddCfgValue("DEFAULT_SERVER", "irc.freenode.net")

	cfg[3].AddCfgValue("köttbulle", "en fin mening!")

	return cfg
}

func TestCfgValue(t *testing.T) {
	cfg := setup()
	cfgvalues := cfg[0].GetCfgValues()
	cfg[0].AddCfgValue("NICK", "Pelle")
	cfg[0].AddCfgValue("DEFAULT_SERVER", "irc.freenode.net")
	if cfgvalues["NICK"] != "Pelle" {
		t.Fail()
	}
	if cfgvalues["DEFAULT_SERVER"] != "irc.freenode.net" {
		t.Fail()
	}

	cfg[0].AddCfgValue("NICK", "Kalle")
	cfg[0].AddCfgValue("DEFAULT_SERVER", "irc.whatever.net")
	if cfgvalues["NICK"] != "Kalle" {
		t.Fail()
	}
	if cfgvalues["DEFAULT_SERVER"] != "irc.whatever.net" {
		t.Fail()
	}

	cfg[0].RemoveCfgValue("NICK")
	cfg[0].RemoveCfgValue("DEFAULT_SERVER")
	if cfgvalues["NICK"] == "Kalle" {
		t.Fail()
	}
	if cfgvalues["DEFAULT_SERVER"] == "irc.whatever.net" {
		t.Fail()
	}
}

func TestLoadAndSave(t *testing.T) {
	var err error
	cfg := setup()

	cfgvalues := cfg[1].GetCfgValues()
	err = cfg[1].Save()
	checkError(err, t)
	err = cfg[1].Load()
	checkError(err, t)

	if len(cfgvalues) != len(cfg[1].GetCfgValues()) {
		t.Fail()
	}
	for k, v := range cfgvalues {
		if v != cfg[1].GetCfgValues()[k] {
			t.Fail()
		}
	}

	cfg = setup()
	err = cfg[2].Save()
	checkError(err, t)
	cfg[2].AddCfgValue("keke", "lul")
	err = cfg[2].Load()
	checkError(err, t)
	if len(cfg[2].GetCfgValues()) != 0 {
		t.Fail()
	}

	cfg[3].Save()
	cfg[3].Load()

	if cfg[3].GetCfgValues()["köttbulle"] != "en fin mening!" {
		t.Fail()
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}
