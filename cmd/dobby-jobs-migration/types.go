package main

import (
	"localdev/dobby-server/internal/pkg/hogwartsforum/parser"
	"localdev/dobby-server/internal/pkg/hogwartsforum/tool"
)

type session struct {
	Tool *tool.Tool
	Conf config
}

type config struct {
	BaseUrl string `json:"baseUrl"`
}

type OldForumUser struct {
	Id         string
	HPHId      string
	Username   string
	Profile    parser.Profile
	TotalItems []parser.ParsedItem
}

type NewForumUser struct {
	Id       string
	Username string
	//InventoryRows *[]smf_shop_inventory_row
	MemberRow    *smf_members_row
	CustomFields *[]smf_themes
	//Items         *[]Item
	NotFoundItems []parser.ParsedItem
}

type smf_members_row struct {
	id_member            string
	member_name          string
	date_registered      string
	posts                string
	id_group             string
	lngfile              string
	last_login           string
	real_name            string
	instant_messages     string
	unread_messages      string
	new_pm               string
	alerts               string
	buddy_list           string
	pm_ignore_list       string
	pm_prefs             string
	mod_prefs            string
	passwd               string
	email_address        string
	personal_text        string
	birthdate            string
	website_title        string
	website_url          string
	show_online          string
	time_format          string
	signature            string
	time_offset          string
	avatar               string
	usertitle            string
	member_ip            string
	member_ip2           string
	secret_question      string
	secret_answer        string
	id_theme             string
	is_activated         string
	validation_code      string
	id_msg_last_visit    string
	additional_groups    string
	smiley_set           string
	id_post_group        string
	total_time_logged_in string
	password_salt        string
	ignore_boards        string
	warning              string
	passwd_flood         string
	pm_receive_from      string
	timezone             string
	tfa_secret           string
	tfa_backup           string
	shopMoney            string
	shopBank             string
	shopInventory_hide   string
	gamesPass            string
	referral             string
	ref_count            string
}

type smf_themes struct {
	id_member string
	id_theme  string
	variable  string
	value     string
}

type smf_shop_inventory_row struct {
	userid    string
	itemid    string
	trading   string
	tradecost string
	date      string
	tradedate string
	fav       string
	name      string
}

type MigratedUser struct {
	Username     string
	OldForumUser *OldForumUser
	NewForumUser *NewForumUser
}
