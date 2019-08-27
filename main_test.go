package main

import (
	"github.com/StarWarsDev/legion-discord-bot/data"
	"testing"
)

func TestLookupUnit(t *testing.T) {
	name := "Darth Vader"
	ldf := "darthvader"

	assert := func(unit *data.Unit) {
		t.Logf("Unit name should be %s and ldf should be %s", name, ldf)
		{
			if unit.LDF != ldf || unit.Name != name {
				t.Errorf("Expected unit %s (%s) but got %s (%s)", name, ldf, unit.Name, unit.LDF)
				t.Fail()
			}
		}
	}

	t.Logf("Given a unit name of %s", name)
	{
		unit := lookupUtil.LookupUnit(name)
		assert(&unit)
	}

	t.Logf("Given a unit LDF of %s", ldf)
	{
		unit := lookupUtil.LookupUnitByLdf(ldf)
		assert(&unit)
	}
}

func TestLookupUpgrade(t *testing.T) {
	name := "Anger"
	ldf := "anger"

	assert := func(upgrade *data.Upgrade) {
		t.Logf("Upgrade name should be %s and ldf should be %s", name, ldf)
		{
			if upgrade.Name != name || upgrade.LDF != ldf {
				t.Errorf("Expected upgrade %s (%s) but got %s (%s)", name, ldf, upgrade.Name, upgrade.LDF)
				t.Fail()
			}
		}
	}

	t.Logf("Given an upgrade name of %s", name)
	{
		upgrade := lookupUtil.LookupUpgrade(name)
		assert(&upgrade)
	}

	t.Logf("Given an upgrade LDF of %s", ldf)
	{
		upgrade := lookupUtil.LookupUpgradeByLdf(ldf)
		assert(&upgrade)
	}
}

func TestLookupCommand(t *testing.T) {
	name := "Master of Evil"
	ldf := "masterofevil"

	assert := func(command *data.CommandCard) {
		t.Logf("Command card name should be %s and ldf should be %s", name, ldf)
		{
			if command.Name != name || command.LDF != ldf {
				t.Errorf("Expected command card %s (%s) but got %s (%s)", name, ldf, command.Name, command.LDF)
				t.Fail()
			}
		}
	}

	t.Logf("Given a command card name of %s", name)
	{
		command := lookupUtil.LookupCommand(name)
		assert(&command)
	}

	t.Logf("Given a command card ldf of %s", ldf)
	{
		command := lookupUtil.LookupCommandCardByLdf(ldf)
		assert(&command)
	}
}
