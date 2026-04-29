package subjects

type Subject string

// NATS subjects
const (
	MatchJoin     Subject = "match.join"
	MatchLeave    Subject = "match.leave"
	GameProvision Subject = "game.provision"
)

type UserItem struct {
	ID   string
	Name string
}

type GameProvisionPayload struct {
	Users []UserItem
}
