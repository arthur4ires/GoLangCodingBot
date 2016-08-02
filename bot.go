package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"io"
)

type User struct{
	Name string
	Id int
}
type jsonConfig struct {
	BotToken string
	User []User
}

var (
	id_usuario int
	user_admin string
	config jsonConfig
	botToken string
	reiGel_ado int
	barionix int
)

const (
	Bot_V     = "v0.7.2"
	//botName                 = "@GoLangCodingBot"
	botName             = "GoLangCodingBot"
	urlApiTranslate     = "https://translate.google.com/translate_tts?ie=UTF-8&tl=pt-BR&client=tw-ob&q="
	Rules               = "rules.txt"
	TDV                 = "txt_da_vergonha.txt"
	tabela_user         = "usuarios"
	tabela_banidos      = "usuarios_banidos"
)

/////////////////////////////FUNÇÕES RELACIONADAS AO BOT/SERVIDOR/DATABASE///////////////////////////////////
func iniciaBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	logError(err)

	bot.Debug = false

	return bot, err
}
func iniciaDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	logError(err)

	return db, err
}
func logError(log_error error) {
	if log_error != nil {
		log.Println("[ERROR]", log_error)

	}
}

/////////////////////////////FUNÇÕES REACIONADAS A PERMISSÕES//////////////////////////////////////////////
func permCheck(id_usuario int) string {
	if id_usuario == reiGel_ado {
		return "ReiGel_ado"
	} else if id_usuario == barionix {
		return "Barionix"
	} else {
		return "false"
	}
}
func kickUser(id_usuario int, ChatID int64, bot *tgbotapi.BotAPI) {
	k := tgbotapi.ChatMemberConfig{
		ChatID: ChatID,
		UserID: id_usuario,
	}
	bot.KickChatMember(k)
}
func returnAdmins(ChatID int64, bot *tgbotapi.BotAPI) []tgbotapi.ChatMember {
	k := tgbotapi.ChatConfig{
		ChatID: ChatID,
	}
	structAdmin, err := bot.GetChatAdministrators(k)
	logError(err)

	return structAdmin
}
func txtdavergonhaArquivo(username string, id_usuario int, motivo string) {
	txt_da_vergonha := leitorArquivo(TDV)
	escreve_txt := txt_da_vergonha + "\n-" + username + " - ID:" + strconv.Itoa(id_usuario)
	escreveArquivo(TDV, escreve_txt)
	log.Printf("[-]O usuario " + username + " de ID:" + strconv.Itoa(id_usuario) + " tentou " + motivo + "!")
}

///////////////////////////////////////////BASE DE DADOS////////////////////////////////////////////////////
func rUser(db sql.DB, username string, tabela string) (id_usuario int) { //Para o @panuto que não entende minha pog mais que esta implicito o que a função faz(olha o comando sql)ele vai puxar pelo nome do usuario o id do telegram que esta salvo no db :)
	rows, err := db.Query("SELECT id_usuario FROM " + tabela + " WHERE username = '" + username + "'")
	logError(err)
	for rows.Next() {
		err = rows.Scan(&id_usuario)
		logError(err)
	}
	defer rows.Close()
	return id_usuario
}
func iUser(db sql.DB, username string, id_usuario int, tabela string) { //Aqui já insere o usuario
	tx, err := db.Begin()
	logError(err)
	stmt, err := tx.Prepare("insert into " + tabela + "(username,id_usuario) values(?,?)")
	logError(err)
	_, err = stmt.Exec(username, id_usuario)
	logError(err)
	tx.Commit()
}

//////////////////////////////////////MENSAGENS/////////////////////////////////////////////////////////////
func mandarMensagem(ChatID int64, Mensagem string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(ChatID, Mensagem)
	msg.ParseMode = "html"

	bot.Send(msg)
}
func mandarFoto(ChatID int64, UserName string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewPhotoUpload(ChatID, "downloads/imagem.jpg")
	msg.Caption = "@" + UserName

	bot.Send(msg)
}
func mandaAudio(ChatID int64, UserName string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewAudioUpload(ChatID, "downloads/audio.mp3")
	msg.Title = "@" + UserName
	bot.Send(msg)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////FUNÇÕES RELACIONADAS A COMANDO/////////////////////////////////
func Verifica_Comando(mensagem string) bool {
	if string(mensagem[0]) == "/" {
		return true
	} else {
		return false
	}
	return false
}
func funcRegras(mensagem string, username string, id_usuario int) string {
	novas_regras := strings.Replace(mensagem, "/func_regras", "", -1) //Tira o comando da mensagem  :D
	user_admin = permCheck(id_usuario)
	if user_admin == "false" {
		txtdavergonhaArquivo(username, id_usuario, "alterar as regras")
		return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n"
	} else {
		escreveArquivo(Rules, novas_regras)
		log.Printf("[-]Regras atualizadas por " + user_admin)
		return "<b>As regras foram atualizadas com sucesso pelo admin " + user_admin + "!</b>"
	}
}
func regrasPrint() string {
	regras := leitorArquivo(Rules)
	return regras + ""
}
func help() string {
	return "<b> Comandos(Usuarios):</b>\n\n/help\n/admins\n/regras\n/txt_da_vergonha\n/imagem\n/tss\n\n<b> Comandos(Admins):</b>\n\n/func_regras\n/kick\n/ban\n/clear\n"
}
func adminsComando(ChatID int64, bot *tgbotapi.BotAPI) string {
	//admins := returnAdmins(ChatID,bot) na proxima versão eu implemento...
	return "<b>Caso ocorra algum problema ,não fale com sua mãe, fale com :\n\n@ReiGel_ado(Programador) ou @Barionix(Criador)</b>"
}
func kickComando(ChatID int64, mensagem string, id_usuario int, username string, db *sql.DB, bot *tgbotapi.BotAPI) string {
	novo_username := tratamentoString(mensagem, "/kick")
	id_usuario_db := rUser(*db, novo_username, tabela_user)
	user_admin = permCheck(id_usuario)
	if user_admin == "false" {
		txtdavergonhaArquivo(username, id_usuario, " kikar um amiguinho do grupo ")
		return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n"
	} else {
		if id_usuario_db == 0 {
			return "<b>O usuario não consta em nossa base de dados!</b>"
		} else {
			kickUser(id_usuario_db, ChatID, bot)
			log.Printf("[-]O usuario " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi removido pelo admin " + user_admin + "!")
			return "<b>O usuario " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi removido do grupo pelo admin " + user_admin + "!\nE que não volte mais ;-;!</b>"
		}
	}
}
func banComando(ChatID int64, mensagem string, id_usuario int, username string, db *sql.DB, bot *tgbotapi.BotAPI) string {
	novo_username := tratamentoString(mensagem, "/ban") //
	id_usuario_db := rUser(*db, novo_username, tabela_user)
	user_admin = permCheck(id_usuario)
	if user_admin == "false" {
		txtdavergonhaArquivo(username, id_usuario, " banir um amiguinho do grupo ")
		return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n"
	} else {
		if id_usuario_db == 0 {
			return "<b>O usuario não consta em nossa base de dados!</b>"
		} else {
			kickUser(id_usuario_db, ChatID, bot)
			iUser(*db, novo_username, id_usuario_db, tabela_banidos)
			log.Printf("[-]O usuarios " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi banido pelo admin " + user_admin + "!")
			return "<b>O usuario " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + "foi banido do grupo pelo admin " + user_admin + "\nEsse não vai voltar mais ;-;!</b>"
		}
	}
}
func tdv() string {
	tdv := leitorArquivo(TDV)
	return "<b>###########MURAL DA VERGONHA###########</b>\n\n" + tdv + ""
}
func imagemD(ChatID int64, username string, mensagem string, bot *tgbotapi.BotAPI, logApi tgbotapi.APIResponse) string {
	url_escaped := tratamentoString(mensagem, "/imagem")
	url_escaped = tratamentoString(url_escaped, botName)
	if url_escaped == "" {
		return "<b>Usage:/imagem http://site.com/image.jpg\nSo vou ensinar uma vez amiguinho 0-0!\n\n"
	}
	if valida := validaUrl(url_escaped); valida == false {
		return "<b>Url invalida!</b>"
	}
	resp := baixarArquivo(url_escaped)
	escreveArquivo("downloads/imagem.jpg", resp)
	mandarFoto(ChatID, username, bot)
	return ""
}
func Clear(mensagem string, id_usuario int, username string) string {
	mensagem = tratamentoString(mensagem, "/clear")
	if permCheck(id_usuario) == "false" {
		txtdavergonhaArquivo(username, id_usuario, " limpar um arquivo ")
		return "<b>Te peguei com a boca na butija ne?\nParabens!\nVocê estava tentando limpar seu nome do Muralzinho...MAIS NAÃO VAI!</b>"
	} else {
		err := os.Remove(mensagem)
		if err != nil {
			logError(err)
			return "<b>O arquivo " + mensagem + " não existe no diretorio local , sé ta bebado né?</b>"
		}
		return "<b>Arquivo " + mensagem + " foi apagado com sucesso!</b>"
	}
}
func ttsTranslate(ChatID int64, mensagem string, username string, bot *tgbotapi.BotAPI) string {
	mensagem = tratamentoString(mensagem, "/tts")
	audio := baixarArquivo(urlApiTranslate + mensagem)
	escreveArquivo("downloads/audio.mp3", audio)
	mandaAudio(ChatID, username, bot)
	return "<b>Audio enviado!</b>"
}
func comandos(mensagem string, id_usuario int, username string, ChatID int64, db *sql.DB, bot *tgbotapi.BotAPI, logApi tgbotapi.APIResponse) string {
	if strings.Contains(mensagem, "/func_regras") {
		return funcRegras(mensagem, username, id_usuario)
	} else if strings.Contains(mensagem, "/regras") {
		return regrasPrint()
	} else if strings.Contains(mensagem, "/help") {
		return help()
	} else if strings.Contains(mensagem, "/admins") {
		return adminsComando(ChatID, bot)
	} else if strings.Contains(mensagem, "/kick") {
		return kickComando(ChatID, mensagem, id_usuario, username, db, bot)
	} else if strings.Contains(mensagem, "/txt_da_vergonha") {
		return tdv()
	} else if strings.Contains(mensagem, "/imagem") {
		return imagemD(ChatID, username, mensagem, bot, logApi)
	} else if strings.Contains(mensagem, "/clear") {
		return Clear(mensagem, id_usuario, username)
	} else if strings.Contains(mensagem, "/ban") {
		return banComando(ChatID, mensagem, id_usuario, username, db, bot)
	} else if strings.Contains(mensagem, "/tts") {
		return ttsTranslate(ChatID, mensagem, username, bot)
	}
	return "Comando não encontrado,use o /help para saber os comandos!"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////FUNÇÕES RELACIONADAS AO SISTEMA//////////////////////////////////////
func leitorArquivo(arquivo string) string {
	conteudo_arquivo, err := ioutil.ReadFile(arquivo)
	if err != nil {
		logError(err)
		log.Printf("[-]O arquivo " + arquivo + "não pode ser lido!")
		return ""
	}
	return string(conteudo_arquivo)
}
func escreveArquivo(arquivo string, conteudo string) {
	os.Remove("imagem.jpg")
	conteudo_arquivo := []byte(conteudo)
	err := ioutil.WriteFile(arquivo, conteudo_arquivo, 0777)
	logError(err)
}
func tratamentoString(mensagem string, comando string) string {
	mensagem = strings.Replace(mensagem, comando, "", -1)
	mensagem = strings.Replace(mensagem, " ", "", -1)
	if strings.Contains(mensagem, "@") {
		mensagem = strings.Replace(mensagem, "@", "", -1)
	}
	return mensagem
}
func jsonS(decode *tgbotapi.User) string {
	a := decode
	out, err := json.Marshal(a)
	logError(err)
	return string(out)
}

//////////////////////////////////////////"NAVEGADOR"////////////////////////////////////////////////////
func baixarArquivo(url_download string) string {
	resp, err := http.Get(url_download)
	logError(err)
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	logError(err)
	return string(b)
}
func validaUrl(url_download string) bool {
	if url_download == "" {
		return false
	}
	url_v, err := url.Parse(url_download)
	if err != nil {
		return false
	}
	if url_v.Host == "" {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	jsonConfigF := leitorArquivo("config/config.json")
	dec := json.NewDecoder(strings.NewReader(jsonConfigF))
	for {
		if err := dec.Decode(&config); err == io.EOF {
			break
		} else if err != nil {
			logError(err)
		}
	}

	botToken = config.BotToken
	reiGel_ado = config.User[0].Id
	barionix = config.User[1].Id

	bot, err := iniciaBot()
	logError(err)
	db, err := iniciaDatabase()
	logError(err)

	defer db.Close()

	log.Printf("Logado na conta %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5 //Tempo pra atualizar
	logErroRequests := tgbotapi.APIResponse{}

	updates, err := bot.GetUpdatesChan(u)

	for msg := range updates { //Recebe as mensagem
		if msg.Message == nil {
			continue
		}
		//////////////////////////////Verifica o BAN//////////////////
		if msg.Message.NewChatMember != nil {
			usuario := fmt.Sprintf("%s", msg.Message.NewChatMember)
			if rUser(*db, string(usuario), tabela_banidos) == 0 {
			} else {
				log.Println("[-]O usuarios " + usuario + " tentou entrar no grupo!")
				kickUser(rUser(*db, string(usuario), tabela_user), msg.Message.Chat.ID, bot)
				msg.Message.NewChatMember = nil
			}
		}
		///////////////////////////////Cadastra o Usuario//////////////////////////////////////
		if rUser(*db, msg.Message.From.UserName, tabela_user) != msg.Message.From.ID {
			iUser(*db, msg.Message.From.UserName, msg.Message.From.ID, tabela_user)
			log.Println("[+]O usuario " + msg.Message.From.UserName + "(ID:" + strconv.Itoa(msg.Message.From.ID) + ") foi cadastrado!")
		}
		/////////////////////////////////////////////////////////////////////////////////////////
		if msg.Message.Text != "" {
			log.Printf("[%s] %s", msg.Message.From.UserName, msg.Message.Text)
		}
		if msg.Message.LeftChatMember != nil {
			log.Printf("[%s] foi removido do grupo.", msg.Message.LeftChatMember)
			mandarMensagem(msg.Message.Chat.ID, "<b>Esse deve ter feito merda....( ͡° ͜ʖ ͡°)</b>", bot)
		}
		if msg.Message.NewChatMember != nil {
			log.Printf("[%s] foi adicionado/convidado ao grupo por %s!", msg.Message.NewChatMember, msg.Message.From.UserName)
			mandarMensagem(msg.Message.Chat.ID, "<b>Eai GOleiro , seja bem vindo a alcateia! Mas conta ai , como chegou aqui ?\n\nVersão do Bot:"+Bot_V+"</b>", bot)
		}

		if msg.Message.Text == "" {
			continue
		} else {
			if Verifica_Comando(msg.Message.Text) == true {
				Resp := comandos(msg.Message.Text, msg.Message.From.ID, msg.Message.From.UserName, msg.Message.Chat.ID, db, bot, logErroRequests) + "\n\n<b>Versão do Bot:" + Bot_V + "</b>"
				mandarMensagem(msg.Message.Chat.ID, Resp, bot)
			}
		}
	}
}
