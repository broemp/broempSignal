package backend

import (
	"fmt"

	"github.com/broemp/broempsignal/model"
	"github.com/bwmarrin/discordgo"
)

type assembleCall struct {
	friends            []model.User
	accepted           []model.User
	declined           []model.User
	discordInteraction *discordgo.Interaction
}

var (
	// map[host discord id] slice of friends
	callList = make(map[string]assembleCall)
	// map[friend telegram id]host discord id
	// check who invited invitee
	inviteeList = make(map[int64]string)
)

// Send messages to all friends
func assemble(userId string, interactionId *discordgo.Interaction) string {
	var user model.User
	// Load user with friends
	db.Preload("Friends").First(&user, "discord_id = ?", userId)

	if len(user.Friends) == 0 {
		return "You have no friends! :("
	}

	//Create assemble call and add to list
	currentAssembleCall := assembleCall{friends: user.Friends, discordInteraction: interactionId}
	callList[userId] = currentAssembleCall

	// Send message to all friends
	for _, f := range currentAssembleCall.friends {
		if Debug {
			println("Sending message to " + f.Name)
		}
		inviteeList[*f.TelegramId] = userId
		if f.TelegramId != nil {
			sendMessage(*f.TelegramId, user.Name)
		}
	}
	return "Sent message to " + fmt.Sprint(len(user.Friends)) + " friends!"
}

// Send message if somebody accpted
func acceptInvite(friend int64) {
	host, ok := inviteeList[friend]
	if !ok {
		return
	}
	delete(inviteeList, friend)

	assembleCall, ok := callList[host]
	if !ok {
		return
	}
	var friend_user model.User
	for i, f := range assembleCall.friends {
		if *f.TelegramId == friend {
			friend_user = f
			assembleCall.friends = removeElement(assembleCall.friends, i)
			assembleCall.accepted = append(assembleCall.accepted, f)
			break
		}
	}

	callList[host] = assembleCall
	updateInteraction(assembleCall.discordInteraction, friend_user.Name+" accepted the invite!")
}

func declineInvite(friend int64) {
	host, ok := inviteeList[friend]
	if !ok {
		return
	}
	delete(inviteeList, friend)

	assembleCall, ok := callList[host]
	if !ok {
		return
	}

	for i, f := range callList[host].friends {
		if f.TelegramId == &friend {
			assembleCall.friends = removeElement(assembleCall.friends, i)
			assembleCall.declined = append(assembleCall.declined, f)
			break
		}
	}

	callList[host] = assembleCall
}

func removeElement(s []model.User, i int) []model.User {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
