package custom_command

import (
	"fmt"
	"strings"

	"github.com/janitorjeff/jeff-bot/core"
	"github.com/janitorjeff/jeff-bot/frontends/discord"

	dg "github.com/bwmarrin/discordgo"
)

var Advanced = advanced{}

type advanced struct{}

func (advanced) Type() core.CommandType {
	return core.Advanced
}

func (advanced) Permitted(m *core.Message) bool {
	return m.Author.Mod()
}

func (advanced) Names() []string {
	return []string{
		"command",
		"cmd",
	}
}

func (advanced) Description() string {
	return "Add, edit, delete or list custom commands."
}

func (c advanced) UsageArgs() string {
	return c.Children().Usage()
}

func (advanced) Category() core.CommandCategory {
	return core.CommandCategoryModerators
}

func (advanced) Examples() []string {
	return nil
}

func (advanced) Parent() core.CommandStatic {
	return nil
}

func (advanced) Children() core.CommandsStatic {
	return core.CommandsStatic{
		AdvancedAdd,
		AdvancedEdit,
		AdvancedDelete,
		AdvancedList,
		AdvancedHistory,
	}
}

func (c advanced) Init() error {
	core.Hooks.Register(c.writeCustomCommand)
	return nil
}

func (advanced) writeCustomCommand(m *core.Message) {
	fields := m.Fields()

	if len(fields) > 1 {
		return
	}

	here, err := m.Here.ScopeLogical()
	if err != nil {
		return
	}

	resp, err := Show(here, fields[0])
	if err != nil {
		return
	}

	m.Write(resp, nil)
}

func (advanced) Run(m *core.Message) (any, error, error) {
	return m.Usage(), core.ErrMissingArgs, nil
}

/////////
//     //
// add //
//     //
/////////

var AdvancedAdd = advancedAdd{}

type advancedAdd struct{}

func (c advancedAdd) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedAdd) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedAdd) Names() []string {
	return core.AliasesAdd
}

func (advancedAdd) Description() string {
	return "Add a command."
}

func (advancedAdd) UsageArgs() string {
	return "<trigger> <text>"
}

func (c advancedAdd) Category() core.CommandCategory {
	return c.Parent().Category()
}

func (advancedAdd) Examples() []string {
	return nil
}

func (advancedAdd) Parent() core.CommandStatic {
	return Advanced
}

func (advancedAdd) Children() core.CommandsStatic {
	return nil
}

func (advancedAdd) Init() error {
	return nil
}

func (c advancedAdd) Run(m *core.Message) (any, error, error) {
	if len(m.Command.Args) < 2 {
		return m.Usage(), core.ErrMissingArgs, nil
	}

	switch m.Frontend.Type() {
	case discord.Frontend.Type():
		return c.discord(m)
	default:
		return c.text(m)
	}
}

func (c advancedAdd) discord(m *core.Message) (*dg.MessageEmbed, error, error) {
	trigger, usrErr, err := c.core(m)
	if err != nil {
		return nil, usrErr, err
	}

	trigger = discord.PlaceInBackticks(trigger)

	embed := &dg.MessageEmbed{
		Description: c.err(usrErr, trigger),
	}

	return embed, usrErr, nil
}

func (c advancedAdd) text(m *core.Message) (string, error, error) {
	trigger, usrErr, err := c.core(m)
	if err != nil {
		return "", usrErr, err
	}

	trigger = fmt.Sprintf("'%s'", trigger)

	return c.err(usrErr, trigger), usrErr, nil
}

func (advancedAdd) err(usrErr error, trigger string) string {
	switch usrErr {
	case nil:
		return fmt.Sprintf("Custom command %s has been added.", trigger)
	case ErrTriggerExists:
		return fmt.Sprintf("Custom command %s already exists.", trigger)
	case ErrBuiltinCommand:
		return fmt.Sprintf("Command %s already exists as a built-in command.", trigger)
	default:
		return "Something went wrong..."
	}
}

func (c advancedAdd) core(m *core.Message) (string, error, error) {
	trigger := m.Command.Args[0]
	response := m.RawArgs(1)

	author, err := m.Author.Scope()
	if err != nil {
		return "", nil, err
	}

	here, err := m.Here.ScopeLogical()
	if err != nil {
		return "", nil, err
	}

	usrErr, err := Add(here, author, trigger, response)
	return trigger, usrErr, err
}

//////////
//      //
// edit //
//      //
//////////

var AdvancedEdit = advancedEdit{}

type advancedEdit struct{}

func (c advancedEdit) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedEdit) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedEdit) Names() []string {
	return core.AliasesEdit
}

func (advancedEdit) Description() string {
	return "Edit a command."
}

func (advancedEdit) UsageArgs() string {
	return "<trigger> <text>"
}

func (c advancedEdit) Category() core.CommandCategory {
	return c.Parent().Category()
}

func (advancedEdit) Examples() []string {
	return nil
}

func (advancedEdit) Parent() core.CommandStatic {
	return Advanced
}

func (advancedEdit) Children() core.CommandsStatic {
	return nil
}

func (advancedEdit) Init() error {
	return nil
}

func (c advancedEdit) Run(m *core.Message) (any, error, error) {
	if len(m.Command.Args) < 2 {
		return m.Usage(), core.ErrMissingArgs, nil
	}

	switch m.Frontend.Type() {
	case discord.Frontend.Type():
		return c.discord(m)
	default:
		return c.text(m)
	}
}

func (c advancedEdit) discord(m *core.Message) (*dg.MessageEmbed, error, error) {
	trigger, usrErr, err := c.core(m)
	if err != nil {
		return nil, usrErr, err
	}

	trigger = discord.PlaceInBackticks(trigger)

	embed := &dg.MessageEmbed{
		Description: c.err(usrErr, trigger),
	}

	return embed, usrErr, nil
}

func (c advancedEdit) text(m *core.Message) (string, error, error) {
	trigger, usrErr, err := c.core(m)
	if err != nil {
		return "", usrErr, err
	}

	trigger = fmt.Sprintf("'%s'", trigger)

	return c.err(usrErr, trigger), usrErr, nil
}

func (advancedEdit) err(usrErr error, trigger string) string {
	switch usrErr {
	case nil:
		return fmt.Sprintf("Custom command %s has been modified.", trigger)
	case ErrTriggerNotFound:
		return fmt.Sprintf("Custom command %s doesn't exist.", trigger)
	default:
		return "Something went wrong..."
	}
}

func (advancedEdit) core(m *core.Message) (string, error, error) {
	trigger := m.Command.Args[0]
	response := m.RawArgs(1)

	author, err := m.Author.Scope()
	if err != nil {
		return "", nil, err
	}

	here, err := m.Here.ScopeLogical()
	if err != nil {
		return "", nil, err
	}

	usrErr, err := Edit(here, author, trigger, response)
	return trigger, usrErr, err
}

////////////
//        //
// delete //
//        //
////////////

var AdvancedDelete = advancedDelete{}

type advancedDelete struct{}

func (c advancedDelete) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedDelete) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedDelete) Names() []string {
	return core.AliasesDelete
}

func (advancedDelete) Description() string {
	return "Delete a command."
}

func (advancedDelete) UsageArgs() string {
	return "<trigger>"
}

func (c advancedDelete) Category() core.CommandCategory {
	return c.Parent().Category()
}

func (advancedDelete) Examples() []string {
	return nil
}

func (advancedDelete) Parent() core.CommandStatic {
	return Advanced
}

func (advancedDelete) Children() core.CommandsStatic {
	return nil
}

func (advancedDelete) Init() error {
	return nil
}

func (c advancedDelete) Run(m *core.Message) (any, error, error) {
	if len(m.Command.Args) < 1 {
		return m.Usage(), core.ErrMissingArgs, nil
	}

	switch m.Frontend.Type() {
	case discord.Frontend.Type():
		return c.discord(m)
	default:
		return c.text(m)
	}
}

func (c advancedDelete) discord(m *core.Message) (*dg.MessageEmbed, error, error) {
	trigger, usrErr, err := c.core(m)
	if err != nil {
		return nil, usrErr, err
	}

	trigger = discord.PlaceInBackticks(trigger)

	embed := &dg.MessageEmbed{
		Description: c.err(usrErr, trigger),
	}

	return embed, usrErr, nil
}

func (c advancedDelete) text(m *core.Message) (string, error, error) {
	trigger, usrErr, err := c.core(m)
	if err != nil {
		return "", usrErr, err
	}

	trigger = fmt.Sprintf("'%s'", trigger)

	return c.err(usrErr, trigger), usrErr, nil
}

func (advancedDelete) err(usrErr error, trigger string) string {
	switch usrErr {
	case nil:
		return fmt.Sprintf("Custom command %s has been deleted.", trigger)
	case ErrTriggerNotFound:
		return fmt.Sprintf("Custom command %s doesn't exist.", trigger)
	default:
		return "Something went wrong..."
	}
}

func (advancedDelete) core(m *core.Message) (string, error, error) {
	trigger := m.Command.Args[0]

	here, err := m.Here.ScopeLogical()
	if err != nil {
		return "", nil, err
	}

	author, err := m.Author.Scope()
	if err != nil {
		return "", nil, err
	}

	usrErr, err := Delete(here, author, trigger)
	return trigger, usrErr, err
}

//////////
//      //
// list //
//      //
//////////

var AdvancedList = advancedList{}

type advancedList struct{}

func (c advancedList) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedList) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedList) Names() []string {
	return core.AliasesList
}

func (advancedList) Description() string {
	return "List commands."
}

func (advancedList) UsageArgs() string {
	return ""
}

func (c advancedList) Category() core.CommandCategory {
	return c.Parent().Category()
}

func (advancedList) Examples() []string {
	return nil
}

func (advancedList) Parent() core.CommandStatic {
	return Advanced
}

func (advancedList) Children() core.CommandsStatic {
	return nil
}

func (advancedList) Init() error {
	return nil
}

func (c advancedList) Run(m *core.Message) (any, error, error) {
	switch m.Frontend.Type() {
	case discord.Frontend.Type():
		return c.discord(m)
	default:
		return c.text(m)
	}
}

func (c advancedList) discord(m *core.Message) (*dg.MessageEmbed, error, error) {
	triggers, err := c.core(m)
	if err != nil {
		return nil, nil, err
	}

	var reply string

	if len(triggers) == 0 {
		reply = "There are no custom commands."
	} else {
		for i := range triggers {
			triggers[i] = "- " + discord.PlaceInBackticks(triggers[i])
		}
		reply = strings.Join(triggers, "\n")
	}

	embed := &dg.MessageEmbed{
		Description: reply,
	}

	return embed, nil, nil
}

func (c advancedList) text(m *core.Message) (string, error, error) {
	triggers, err := c.core(m)
	if err != nil {
		return "", nil, err
	}

	if len(triggers) == 0 {
		return "There are no custom commands.", nil, nil
	}
	return strings.Join(triggers, ", "), nil, nil
}

func (c advancedList) core(m *core.Message) ([]string, error) {
	here, err := m.Here.ScopeLogical()
	if err != nil {
		return nil, err
	}
	return List(here)
}

/////////////
//         //
// history //
//         //
/////////////

var AdvancedHistory = advancedHistory{}

type advancedHistory struct{}

func (c advancedHistory) Type() core.CommandType {
	return c.Parent().Type()
}

func (c advancedHistory) Permitted(m *core.Message) bool {
	return c.Parent().Permitted(m)
}

func (advancedHistory) Names() []string {
	return []string{
		"history",
	}
}

func (advancedHistory) Description() string {
	return "View a command's entire history of changes."
}

func (advancedHistory) UsageArgs() string {
	return "<trigger>"
}

func (c advancedHistory) Category() core.CommandCategory {
	return c.Parent().Category()
}

func (advancedHistory) Examples() []string {
	return nil
}

func (advancedHistory) Parent() core.CommandStatic {
	return Advanced
}

func (advancedHistory) Children() core.CommandsStatic {
	return nil
}

func (advancedHistory) Init() error {
	return nil
}

func (c advancedHistory) Run(m *core.Message) (any, error, error) {
	if len(m.Command.Args) < 1 {
		return m.Usage(), core.ErrMissingArgs, nil
	}

	switch m.Frontend.Type() {
	case discord.Frontend.Type():
		return c.discord(m)
	default:
		return nil, nil, nil
	}
}

func formatTime(timestamp int64) string {
	return fmt.Sprintf("<t:%d:D>", timestamp)
}

func formatCreate(timestamp int64, response string) string {
	when := formatTime(timestamp)
	return fmt.Sprintf("created '%s' %s by @", response, when)
}

func formatModify(timestamp int64, response string) string {
	when := formatTime(timestamp)
	return fmt.Sprintf("modified to '%s' %s by @", response, when)
}

func formatDelete(timestamp int64) string {
	when := formatTime(timestamp)
	return fmt.Sprintf("deleted %s by @", when)
}

func (c advancedHistory) discord(m *core.Message) (*dg.MessageEmbed, error, error) {
	trigger, history, err := c.core(m)
	if err != nil {
		return nil, nil, err
	}

	if len(history) == 0 {
	}

	const zeroWidthSpace = "\u200b"

	var action []string
	var response []string
	var when []string

	for i := 0; i < len(history); i++ {
		hist := history[i]

		if i == 0 {
			// creation
			action = append(action, "created")
			response = append(response, hist.response)
			when = append(when, formatTime(hist.created))
		} else if history[i-1].deleted == hist.created {
			// modification
			action = append(action, "edited")
			response = append(response, hist.response)
			when = append(when, formatTime(hist.created))
		} else {
			// deletion
			action = append(action, "deleted")
			response = append(response, "")
			when = append(when, formatTime(history[i-1].deleted))

			action = append(action, "created")
			response = append(response, hist.response)
			when = append(when, formatTime(hist.created))
		}

		if i == len(history)-1 && hist.deleted != 0 {
			action = append(action, "deleted")
			response = append(response, "")
			when = append(when, formatTime(hist.deleted))
		}
	}

	embed := &dg.MessageEmbed{
		Title: discord.PlaceInBackticks(trigger),
		Fields: []*dg.MessageEmbedField{
			{
				Name:   "action",
				Value:  strings.Join(action, "\n"),
				Inline: true,
			},
			{
				Name:   "response",
				Value:  strings.Join(response, "\n"),
				Inline: true,
			},
			{
				Name:   "when",
				Value:  strings.Join(when, "\n"),
				Inline: true,
			},
		},
	}

	return embed, nil, nil
}

func (advancedHistory) core(m *core.Message) (string, []customCommand, error) {
	trigger := m.Command.Args[0]

	here, err := m.Here.ScopeLogical()
	if err != nil {
		return trigger, nil, err
	}

	history, err := History(here, trigger)
	return trigger, history, err
}
