package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/telegram-bot-api.v4"
	//"fmt"
)

var (
	ReiGel_ado int = 0
	Barionix   int = 0 //Gay
	id_usuario int
	user_admin string
)

const (
	//Bot_Token = "Malandro" //Meu
	Bot_Token = "Malandro" //GoLafTest
	Bot_V     = " v0.5"
	Bot_Name  = "@GoLangCodingBot"
	Rules     = "rules.txt"
	TDV       = "txt_da_vergonha.txt"
)

func error_check(log_error error) {
	if log_error != nil {
		log.Fatal(log_error)
	}
}

/////////////////////////////FUNÇÕES REACIONADAS A PERMISSÕES//////////////////////////////////////////////
func Permissao(id_usuario int) string {
	if id_usuario == ReiGel_ado {
		return "ReiGel_ado"
	} else if id_usuario == Barionix {
		return "Barionix"
	} else {
		return "Hasker"
	}
}
func Kick(id_usuario int, ChatID int64, bot *tgbotapi.BotAPI) {
	k := tgbotapi.ChatMemberConfig{
		ChatID: ChatID,
		UserID: id_usuario,
	}
	bot.KickChatMember(k)
}

/////////////////////////////FUNÇÕES RELACIONADAS AO BOT/SERVIDOR/DATABASE///////////////////////////////////
func Inicia_Bot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(Bot_Token)
	error_check(err)

	bot.Debug = false

	return bot, err
}
func Inicia_Database() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	error_check(err)

	return db, err
}

///////////////////////////////////////////BASE DE DADOS////////////////////////////////////////////////////
func Retorna_User(db *sql.DB, username string) (id_usuario int) {
	rows, err := db.Query("SELECT id_usuario FROM usuarios WHERE username = '" + username + "'")
	error_check(err)
	for rows.Next() {
		err = rows.Scan(&id_usuario)
		error_check(err)
	}
	defer rows.Close()
	return id_usuario
}
func Insere_User(db *sql.DB, username string, id_usuario int) {
	tx, err := db.Begin()
	error_check(err)
	stmt, err := tx.Prepare("insert into usuarios(username,id_usuario) values(?,?)")
	error_check(err)
	_, err = stmt.Exec(username, id_usuario)
	error_check(err)
	tx.Commit()
}

//////////////////////////////////////MENSAGENS/////////////////////////////////////////////////////////////
func Mandar_Mensagen(ChatID int64, Mensagem string) {
	bot, err := Inicia_Bot()
	error_check(err)

	msg := tgbotapi.NewMessage(ChatID, Mensagem)
	msg.ParseMode = "html"

	bot.Send(msg)
}

func Responder_Mensagens(ChatID int64, Mensagem string, MensagemID int) {
	bot, err := Inicia_Bot()
	error_check(err)

	msg := tgbotapi.NewMessage(ChatID, Mensagem)
	msg.ReplyToMessageID = MensagemID
	msg.ParseMode = "html"

	bot.Send(msg)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////FUNÇÕES RELACIONADAS A COMANDO/////////////////////////////////
func Verifica_Comando(mensagem string) bool {
	if string(mensagem[0]) == "/" {
		return true
	}
	return false
}
func Comandos(mensagem string, id_usuario int, username string, ChatID int64, db *sql.DB, bot *tgbotapi.BotAPI) string {
	///////////////////////REGRAS//////////////////////////////////QQQQ
	if strings.Contains(mensagem, "/func_regras") {
		novas_regras := strings.Replace(mensagem, "/func_regras", "", -1) //Tira o comando da mensagem  :D
		user_admin = Permissao(id_usuario)
		if user_admin == "Hasker" {
			txt_da_vergonha := Leitor(TDV)
			escreve_txt := txt_da_vergonha + "\n-" + username + " - ID:" + strconv.Itoa(id_usuario)
			Escreve(TDV, escreve_txt)
			log.Printf("[-]O usuario " + username + " de ID:" + strconv.Itoa(id_usuario) + " tentou alterar as regras!")
			return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n\n\n<b>Versão do Bot:" + Bot_V + "</b>"
		} else {
			Escreve(Rules, novas_regras)
			log.Printf("[-]Regras atualizadas por " + user_admin)
			return "<b>As regras foram atualizadas com sucesso pelo admin " + user_admin + "!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
		}
	} else if strings.Contains(mensagem, "/regras") {
		regras := Leitor(Rules)
		return regras + "\n\n<b>Versão do Bot:" + Bot_V + "</b>"
		////////////////////////////////////////////////////////////////////Q
	} else if strings.Contains(mensagem, "/help") {
		return "<b>Comandos(Usuarios):</b>\n\n/help\n/admins\n/regras\n/txt_da_vergonha\n\n<b>Comandos(Admins):</b>\n\n/func_regras - Editas as regras!\n\n<b>Versão do Bot:" + Bot_V + "</b>"
	} else if strings.Contains(mensagem, "/admins") {
		return "<b>Caso ocorra algum problema ,não fale com sua mãe, fale com :</b>\n\n@ReiGel_ado<b> ou </b>@Barionix\n\n<b>Versão do Bot:" + Bot_V + "</b>"
	} else if strings.Contains(mensagem, "/kick") {
		novo_username := strings.Replace(mensagem, "/kick", "", -1)
		novo_username = strings.Replace(novo_username, " ", "", -1)
		id_usuario_db := Retorna_User(db, novo_username)
		user_admin = Permissao(id_usuario)
		if user_admin == "Hasker" {
			txt_da_vergonha := Leitor(TDV)
			escreve_txt := txt_da_vergonha + "\n-" + username + " - ID:" + strconv.Itoa(id_usuario)
			Escreve(TDV, escreve_txt)
			log.Printf("[-]O usuario " + username + " de ID:" + strconv.Itoa(id_usuario) + " tentou kikar um membro!")
			return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n\n\n<b>Versão do Bot:" + Bot_V + "</b>"
		} else {
			Kick(id_usuario_db, ChatID, bot)
			if id_usuario_db == 0 {
				return "<b>O usuario não consta em nossa base de dados!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
			} else {
				log.Printf("[-]O usuario " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi removido pelo admin " + user_admin + "!")
				return "<b>O usuario " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi removido do grupo pelo admin " + user_admin + "!\nE que não volte mais ;-;!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
			}
		}
	} else if strings.Contains(mensagem, "/txt_da_vergonha") {
		tdv := Leitor(TDV)
		return "<b>###########MURAL DA VERGONHA###########</b>\n\n" + tdv + "\n\n<b>Versão do Bot:" + Bot_V + "</b>"
	}
	return "Comando não encontrado,use o /help para saber os comandos!"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////FUNÇÕES RELACIONADAS AO SISTEMA//////////////////////////////////////
func Leitor(arquivo string) string {
	conteudo_arquivo, err := ioutil.ReadFile(arquivo)
	error_check(err)
	return string(conteudo_arquivo)
}
func Escreve(arquivo string, conteudo string) {
	conteudo_arquivo := []byte(conteudo)
	err := ioutil.WriteFile(arquivo, conteudo_arquivo, 0777)
	error_check(err)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	bot, err := Inicia_Bot()
	error_check(err)
	db, err := Inicia_Database()
	error_check(err)

	defer db.Close()

	log.Printf("Logado na conta %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0 //Tempo pra atualizar

	updates, err := bot.GetUpdatesChan(u)

	for msg := range updates { //Recebe as mensagem
		if msg.Message == nil {
			continue
		}
		if Retorna_User(db, msg.Message.From.UserName) == msg.Message.From.ID {
		} else {
			Insere_User(db, msg.Message.From.UserName, msg.Message.From.ID)
			log.Println("O usuario " + msg.Message.From.UserName + "(ID:" + strconv.Itoa(msg.Message.From.ID) + ") foi cadastrado!")
		}

		log.Printf("[%s] %s", msg.Message.From.UserName, msg.Message.Text)

		if msg.Message.LeftChatMember != nil { //Usuario Saiu
			log.Printf("[%s] foi removido do grupo.", msg.Message.LeftChatMember)
			Mandar_Mensagen(msg.Message.Chat.ID, "<b>Esse deve ter feito merda....( ͡° ͜ʖ ͡°)</b>\n\n<b>Versão do Bot:"+Bot_V+"</b>")
		}
		if msg.Message.NewChatMember != nil {
			log.Printf("[%s] foi adicionado/convidado ao grupo por %s!", msg.Message.NewChatMember, msg.Message.From.UserName)
			Mandar_Mensagen(msg.Message.Chat.ID, "<b>Eai GOleiro , seja bem vindo a alcateia! Mas conta ai , como chegou aqui ?</b>\n\n<b>Versão do Bot:"+Bot_V+"</b>")
		}

		if msg.Message.Text == "" {
			continue
		} else {
			if Verifica_Comando(msg.Message.Text) == true {
				x := Comandos(msg.Message.Text, msg.Message.From.ID, msg.Message.From.UserName, msg.Message.Chat.ID, db, bot)
				Mandar_Mensagen(msg.Message.Chat.ID, x)
			}
		}
	}
}
