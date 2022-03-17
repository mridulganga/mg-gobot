package handlers

import (
	"fmt"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/mridulganga/mg-gobot/pkg/db"
	"github.com/mridulganga/mg-gobot/pkg/models"
	"github.com/mridulganga/mg-gobot/pkg/utils"
)

type Mono struct {
	db db.DB
}

func NewMono(dbObj db.DB) Mono {
	return Mono{
		db: dbObj,
	}
}

func (h Mono) getOrCreateUser(id, username string) *models.User {
	u := h.db.GetUser(id)
	if u == nil {
		u = &models.User{
			ID:   id,
			Name: username,
			Balance: models.Balance{
				Wallet: 0,
				Bank:   0,
			},
		}
		h.db.PutUser(*u)
	}
	return u
}

func (h Mono) BalanceHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := h.getOrCreateUser(m.Author.ID, m.Author.Username)
	title := fmt.Sprintf("Balance @%v", u.Name)
	msg := fmt.Sprintf("Wallet: %v\nBank: %v", u.Balance.Wallet, u.Balance.Bank)
	utils.SendEmbed(s, m, title, msg)
}

func (h Mono) RichHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	users := h.db.ListUsers()
	log.Info(users)

	sort.Slice(users, func(i, j int) bool {
		return users[i].Balance.Wallet+users[i].Balance.Bank > users[j].Balance.Wallet+users[j].Balance.Bank
	})
	msg := ""
	if len(users) > 3 {
		users = users[:3]
	}
	for _, u := range users {
		msg += fmt.Sprintf("%v: %v\n", u.Name, u.Balance.Wallet+u.Balance.Bank)
	}
	utils.SendEmbed(s, m, "Rich People", msg)

}

func (h Mono) SendHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	fromUser := h.getOrCreateUser(m.Author.ID, m.Author.Username)
	if len(m.Mentions) != 1 {
		utils.SendReply(s, m, "Whom do you want to send the money?")
		return
	}
	toUser := h.getOrCreateUser(m.Mentions[0].ID, m.Mentions[0].Username)

	amount, present, err := utils.ParseAmount(m, 3)
	if err != nil {
		utils.SendReply(s, m, "Please enter a valid amount")
		return
	}
	if !present {
		amount = fromUser.Balance.Wallet
	}
	if amount < 1 {
		utils.SendReply(s, m, fmt.Sprintf("You can't send %v", amount))
		return
	}

	if fromUser.Balance.Wallet >= amount {
		fromUser.Balance.Wallet -= amount
		toUser.Balance.Wallet += amount
		h.db.PutUser(*fromUser)
		h.db.PutUser(*toUser)
		utils.SendEmbed(s, m, "Transaction", fmt.Sprintf("%s sent %v to %s", fromUser.Name, amount, toUser.Name))
		return
	}
	utils.SendReply(s, m, fmt.Sprintf("<@%s> you ain't got that kinda money", fromUser.ID))
}

func (h Mono) BegHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := h.getOrCreateUser(m.Author.ID, m.Author.Username)
	if u.Last.Beg.Before(time.Now().Add(-time.Second * 10)) {
		begAmount := utils.RandomNumberBetween(20, 150)
		begLine := utils.RandomChoice(utils.GetData("beg_lines"))
		donor := utils.RandomChoice(utils.GetData("donors"))

		u.Balance.Wallet += begAmount
		u.Last.Beg = time.Now()
		h.db.PutUser(*u)

		title := fmt.Sprintf("%s donated %v", donor, begAmount)
		msg := fmt.Sprintf("<@%s> %s", u.ID, begLine)
		utils.SendEmbed(s, m, title, msg)
		return
	}
	utils.SendReply(s, m, "you're begging too much, stop it!!")
}

func (h Mono) SearchHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := h.getOrCreateUser(m.Author.ID, m.Author.Username)
	if u.Last.Search.Before(time.Now().Add(-time.Second * 15)) {
		searchAmount := utils.RandomNumberBetween(30, 120)
		searchLine := utils.RandomChoice(utils.GetData("search_lines"))

		u.Balance.Wallet += searchAmount
		u.Last.Search = time.Now()
		h.db.PutUser(*u)

		msg := fmt.Sprintf("<@%s> found %v %s", u.ID, searchAmount, searchLine)
		utils.SendEmbed(s, m, "Yay!", msg)
		return
	}
	utils.SendReply(s, m, "you're searching too hard, it's not there!!")
}

func (h Mono) StealHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	thief := h.getOrCreateUser(m.Author.ID, m.Author.Username)
	victim := h.getOrCreateUser(m.Mentions[0].ID, m.Mentions[0].Username)

	if thief.Balance.Wallet >= 200 && victim.Balance.Wallet >= 200 {
		stealSuccess := utils.YesNo()
		msg := ""
		if stealSuccess {
			stealAmount := utils.RandomNumberBetween(0, victim.Balance.Wallet)
			thief.Balance.Wallet += stealAmount
			victim.Balance.Wallet -= stealAmount
			msg = fmt.Sprintf("<@%s> stole %v from <@%s>", thief.ID, stealAmount, victim.ID)
		} else {
			lossAmount := utils.RandomNumberBetween(0, thief.Balance.Wallet)
			thief.Balance.Wallet -= lossAmount
			victim.Balance.Wallet += lossAmount
			msg = fmt.Sprintf("<@%s> got caught stealing from <@%s> and paid them %v", thief.ID, victim.ID, lossAmount)
		}
		h.db.PutUser(*thief)
		h.db.PutUser(*victim)
		utils.SendEmbed(s, m, "Theft Report", msg)
		return
	}
	utils.SendReply(s, m, "Both the thief and victim need to have atleat 200 in their wallet to steal")
}

func (h Mono) DepositHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := h.getOrCreateUser(m.Author.ID, m.Author.Username)
	amount, present, err := utils.ParseAmount(m, 2)
	if err != nil {
		utils.SendReply(s, m, "Please enter a valid amount")
		return
	}
	if !present {
		amount = u.Balance.Wallet
	}
	if amount < 1 {
		utils.SendReply(s, m, fmt.Sprintf("You can't deposit %v", amount))
		return
	}
	if u.Balance.Wallet >= amount {
		u.Balance.Bank += amount
		u.Balance.Wallet -= amount
		h.db.PutUser(*u)
		utils.SendEmbed(s, m, "Bank Deposit", fmt.Sprintf("<@%s> deposited %v in their bank", u.ID, amount))
	} else {
		utils.SendReply(s, m, "You ain't got that kind of money")
	}
}

func (h Mono) WithdrawHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := h.getOrCreateUser(m.Author.ID, m.Author.Username)
	amount, present, err := utils.ParseAmount(m, 2)
	if err != nil {
		utils.SendReply(s, m, "Please enter a valid amount")
		return
	}
	if !present {
		amount = u.Balance.Bank
	}
	if amount < 1 {
		utils.SendReply(s, m, fmt.Sprintf("You can't withdraw %v", amount))
		return
	}
	if u.Balance.Bank >= amount {
		u.Balance.Wallet += amount
		u.Balance.Bank -= amount
		h.db.PutUser(*u)
		utils.SendEmbed(s, m, "Bank Withdrawal", fmt.Sprintf("<@%s> withdrew %v from their bank", u.ID, amount))
	} else {
		utils.SendReply(s, m, "You ain't got that kind of money")
	}
}

func (h Mono) GambleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	u := h.getOrCreateUser(m.Author.ID, m.Author.Username)

	amount, present, err := utils.ParseAmount(m, 2)
	if err != nil {
		utils.SendReply(s, m, "Please enter a valid amount")
		return
	}
	if !present {
		amount = u.Balance.Wallet
	}

	if amount < 1 {
		utils.SendReply(s, m, fmt.Sprintf("You can't gamble %v", amount))
		return
	}

	u.Balance.Wallet -= amount
	title := ""
	msg := ""
	win, returns := utils.GambleMoney(amount)
	if win {
		title = "Won Gamble"
		msg = fmt.Sprintf("<@%s> won %v while Gambling", u.ID, returns-amount)
	} else {
		title = "Lost Gamble"
		msg = fmt.Sprintf("<@%s> lost %v while Gambling", u.ID, amount-returns)

	}
	u.Balance.Wallet += returns
	h.db.PutUser(*u)
	utils.SendEmbed(s, m, title, msg)
}
