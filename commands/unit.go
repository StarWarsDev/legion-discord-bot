package commands

import (
	"strconv"
	"strings"

	"github.com/StarWarsDev/legion-discord-bot/internal/data"
	"github.com/StarWarsDev/legion-discord-bot/output"
	"github.com/bwmarrin/discordgo"
)

// Unit handles the !unit command
func Unit(unit *data.Unit) discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	costField := output.Field("Cost", strconv.Itoa(unit.Cost))
	rankField := output.Field("Rank", unit.Rank)
	typeField := output.Field("Type", unit.Type)
	factionField := output.Field("Faction", unit.Faction)
	uniqueField := output.Field("Unique", strconv.FormatBool(unit.Unique))

	fields = append(
		fields,
		&costField,
		&rankField,
		&typeField,
		&factionField,
		&uniqueField,
	)

	if unit.Wounds > 0 {
		woundsField := output.Field("Wounds", strconv.Itoa(unit.Wounds))
		fields = append(fields, &woundsField)
	}

	if unit.Courage != nil {
		courageField := output.Field("Courage", strconv.Itoa(*unit.Courage))
		if *unit.Courage < 1 {
			courageField.Value = "-"
		}
		fields = append(fields, &courageField)
	}

	if unit.Resilience != nil {
		resilienceField := output.Field("Resilience", strconv.Itoa(*unit.Resilience))
		if *unit.Resilience < 1 {
			resilienceField.Value = "-"
		}
		fields = append(fields, &resilienceField)
	}

	if unit.Defense != "" {
		defenseField := output.Field("Defense", unit.Defense)
		fields = append(fields, &defenseField)
	}

	if len(unit.Requirements) > 0 {
		requirementsField := output.Field("Requirements", strings.Join(unit.Requirements, ", "))
		fields = append(fields, &requirementsField)
	}

	if len(unit.Keywords) > 0 {
		keywordsField := output.Field("Keywords", strings.Join(unit.Keywords, ", "))
		fields = append(fields, &keywordsField)
	}

	if len(unit.Entourage) > 0 {
		entourageField := output.Field("Entourage", strings.Join(unit.Entourage, ", "))
		fields = append(fields, &entourageField)
	}

	if unit.Surge != nil && unit.Surge.String() != "" {
		surgeField := output.Field("Surge", unit.Surge.String())
		fields = append(fields, &surgeField)
	}

	if len(unit.Weapons) > 0 {
		var weapons []string
		for _, weapon := range unit.Weapons {
			weapons = append(weapons, weapon.String())
		}
		weaponsField := output.Field("Weapons", strings.Join(weapons, "\n\n"))
		fields = append(fields, &weaponsField)
	}

	if len(unit.Slots) > 0 {
		slotsField := output.Field("Slots", strings.Join(unit.Slots, ", "))
		fields = append(fields, &slotsField)
	}

	if len(unit.CommandCards) > 0 {
		var names []string
		for _, commandCard := range unit.CommandCards {
			names = append(names, commandCard.Name)
		}
		commandCardsField := output.Field("Command Cards", strings.Join(names, ", "))
		fields = append(fields, &commandCardsField)
	}

	response := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name: "Unit",
		},
		Color:  0xffffff,
		Title:  unit.Name,
		Fields: fields,
		Image: &discordgo.MessageEmbedImage{
			URL: unit.Image,
		},
	}

	return response
}
