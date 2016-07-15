package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "strings"
    "io/ioutil"
    "strconv"
    "database/sql"
    "os"
    "net/http"
    "encoding/json"
    _ "github.com/mattn/go-sqlite3"
    _"fmt"
    _"io"
)

var(
    id_usuario          int
    user_admin          string 
    banStruct            = new(banMember)

)

const ( 
    Bot_Token = "Maohe" //Meu
    //Bot_Token               = "" //GoLafTest
    Bot_V                   = "v0.6.1"
    Bot_Name                = "@GoLangCodingBot"
    Rules                   = "rules.txt" 
    TDV                     = "txt_da_vergonha.txt"
    tabela_user             = "usuarios"
    tabela_banidos          = "usuarios_banidos"
    ReiGel_ado          int = 0
    Barionix            int = 0
)

type banMember struct {
    Id int `json:"id"`
    First_name string `json:"first_name"`
    Last_name string `json:"last_name"`
    Username string `json:"username"`
}

/////////////////////////////FUNÇÕES RELACIONADAS AO BOT/SERVIDOR/DATABASE///////////////////////////////////
func Inicia_Bot()(*tgbotapi.BotAPI, error){
    bot, err := tgbotapi.NewBotAPI(Bot_Token)
    error_check(err)

    bot.Debug = false

    return bot, err
}
func Inicia_Database()(*sql.DB,error){
    db , err := sql.Open("sqlite3","./database.db")
    error_check(err)

    return db,err
}
func error_check(log_error error){
    if log_error != nil {
        log.Println("[ERROR]",log_error)
    }
}
/////////////////////////////FUNÇÕES REACIONADAS A PERMISSÕES//////////////////////////////////////////////
func Permissao(id_usuario int ) string{
    if id_usuario == ReiGel_ado{
        return "ReiGel_ado"
    }else if id_usuario == Barionix{
        return "Barionix"
    }else{
        return "false"
    }
}
func Kick(id_usuario int,ChatID int64,bot *tgbotapi.BotAPI){
    k := tgbotapi.ChatMemberConfig{
        ChatID: ChatID,
        UserID: id_usuario,
    }
    bot.KickChatMember(k)
}
func Retorna_Admins(ChatID int64,bot *tgbotapi.BotAPI)[]tgbotapi.ChatMember{
    k := tgbotapi.ChatConfig{
        ChatID: ChatID,
    }
    slice_admin , err := bot.GetChatAdministrators(k)
    error_check(err)

    return slice_admin
}
func Txt_Da_Vergonha_Escreve(username string , id_usuario int , motivo string){
    txt_da_vergonha := Leitor(TDV)
    escreve_txt := txt_da_vergonha + "\n-" + username + " - ID:" + strconv.Itoa(id_usuario)
    Escreve(TDV,escreve_txt)
    log.Printf("[-]O usuario " + username + " de ID:" + strconv.Itoa(id_usuario) + " tentou " + motivo + "!")
}
///////////////////////////////////////////BASE DE DADOS////////////////////////////////////////////////////
func rUser(db sql.DB, username string,tabela string)(id_usuario int){
    rows, err := db.Query("SELECT id_usuario FROM " + tabela + " WHERE username = '" + username + "'")
    error_check(err)
    for rows.Next(){
        err = rows.Scan(&id_usuario)
        error_check(err)
    }
    defer rows.Close()
    return id_usuario
}
func iUser(db sql.DB,username string,id_usuario int,tabela string){
    tx, err := db.Begin()
    error_check(err)
    stmt, err := tx.Prepare("insert into " + tabela +"(username,id_usuario) values(?,?)")
    error_check(err)
    _, err = stmt.Exec(username,id_usuario)
    error_check(err)
    tx.Commit()
}
//////////////////////////////////////MENSAGENS/////////////////////////////////////////////////////////////
func Mandar_Mensagen(ChatID int64,Mensagem string){
    bot,err := Inicia_Bot()
    error_check(err)

    msg := tgbotapi.NewMessage(ChatID,Mensagem)
    msg.ParseMode = "html"

    bot.Send(msg)
}
func Mandar_Foto(ChatID int64,UserName string,file string) {
    bot,err := Inicia_Bot()
    error_check(err)

    msg := tgbotapi.NewPhotoUpload(ChatID, "downloads/" + file)
    msg.Caption = "@" + UserName

    bot.Send(msg)
}

func Responder_Mensagens(ChatID int64,Mensagem string,MensagemID int){
    bot,err := Inicia_Bot()
    error_check(err)

    msg := tgbotapi.NewMessage(ChatID,Mensagem)
    msg.ReplyToMessageID = MensagemID
    msg.ParseMode = "html"

    bot.Send(msg)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////FUNÇÕES RELACIONADAS A COMANDO/////////////////////////////////
func Verifica_Comando(mensagem string)bool{
    if string(mensagem[0]) == "/"{
        return true
    }else{
        return false
    }
    return false
}
func Func_Regras(mensagem string,username string,id_usuario int)string{
    novas_regras := strings.Replace(mensagem,"/func_regras","",-1) //Tira o comando da mensagem  :D
    user_admin = Permissao(id_usuario)
    if user_admin == "false"{
        Txt_Da_Vergonha_Escreve(username,id_usuario,"alterar as regras")
        return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }else{
        Escreve(Rules,novas_regras)
        log.Printf("[-]Regras atualizadas por " + user_admin )
        return "<b>As regras foram atualizadas com sucesso pelo admin " + user_admin + "!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }
}
func Regras()string{
    regras := Leitor(Rules)
    return regras + "\n\n<b>Versão do Bot:" + Bot_V + "</b>"
}
func Help()string{
    return "<b>Comandos(Usuarios):</b>\n\n/help\n/admins\n/regras\n/txt_da_vergonha\n/imagem\n\n<b>Comandos(Admins):</b>\n\n/func_regras\n/kick\n/ban\n/clear\n\n\n<b>Versão do Bot:" + Bot_V + "</b>"
}
func Admins_C(ChatID int64,bot *tgbotapi.BotAPI)string{
    //slice_admins := Retorna_Admins(ChatID,bot) Implementação futura,tenho que aprender a trabalhar com a bosta dos types
    return "<b>Caso ocorra algum problema ,não fale com sua mãe, fale com :\n\n@ReiGel_ado(Programador) ou @Barionix(Criador)\n\nVersão do Bot:" + Bot_V + "</b>"
}
func Kick_Comando(ChatID int64,mensagem string,id_usuario int,username string,db *sql.DB,bot *tgbotapi.BotAPI) string{
    novo_username := Tratamento(mensagem,"/kick")
    id_usuario_db := rUser(*db,novo_username,tabela_user)
    user_admin     =  Permissao(id_usuario)
    if user_admin == "false"{
        Txt_Da_Vergonha_Escreve(username,id_usuario," kikar um amiguinho do grupo ")
        return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }else{
        if id_usuario_db == 0{
            return "<b>O usuario não consta em nossa base de dados!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
        }else{
            Kick(id_usuario_db,ChatID,bot)
            log.Printf("[-]O usuario " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi removido pelo admin " + user_admin + "!")
            return "<b>O usuario " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi removido do grupo pelo admin " + user_admin + "!\nE que não volte mais ;-;!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
        }
    }
}
func Ban_Comando(ChatID int64,mensagem string,id_usuario int ,username string,db *sql.DB,bot *tgbotapi.BotAPI)string{
    novo_username := Tratamento(mensagem,"/ban")//
    id_usuario_db := rUser(*db,novo_username,tabela_user)
    user_admin     = Permissao(id_usuario)
    if user_admin == "false"{
        Txt_Da_Vergonha_Escreve(username,id_usuario, " banir um amiguinho do grupo ")
        return "<b>Você esta muito gracioso @" + username + "\nParabens!\nSeu nome esta no txt da vergonha ;-;!</b>\n\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }else{
        if  id_usuario_db == 0{
            return "<b>O usuario não consta em nossa base de dados!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
        }else{
            Kick(id_usuario_db,ChatID,bot)
            iUser(*db,novo_username,id_usuario_db,tabela_banidos)
            log.Printf("[-]O usuarios " + novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + " foi banido pelo admin " + user_admin + "!")
            return "<b>O usuario " +  novo_username + " de ID:" + strconv.Itoa(id_usuario_db) + "foi banido do grupo pelo admin " + user_admin + "\nEsse não vai voltar mais ;-;!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
        }
    }
}
func Txt_Da_Vergonha()string{
    tdv := Leitor(TDV)
    return "<b>###########MURAL DA VERGONHA###########</b>\n\n" + tdv + "\n\n<b>Versão do Bot:" + Bot_V + "</b>"
}
func Imagem_D(ChatID int64,username string,mensagem string)string{
    url_escaped := Tratamento(mensagem,"/imagem")
    resp := Baixar_Arquivo(url_escaped)
    Escreve("downloads/imagem.jpg",resp)
    Mandar_Foto(ChatID,username,"imagem.jpg")
    return "<b>Imagem enviada com sucesso!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
}
func Clear(mensagem string,id_usuario int,username string)string{
    mensagem = Tratamento(mensagem,"/clear")
    if Permissao(id_usuario) == "false"{
        Txt_Da_Vergonha_Escreve(username,id_usuario," limpar um arquivo ")
        return "<b>Te peguei com a boca na butija ne?\nParabens!\nVocê estava tentando limpar seu nome do Muralzinho...MAIS NAÃO VAI!</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>"
    }else{
        err := os.Remove(mensagem)
        if err != nil{
            error_check(err)
            return "<b>O arquivo " + mensagem + " não existe no diretorio local , sé ta bebado né?</b>"
        }
        return "<b>Arquivo " + mensagem + " foi apagado com sucesso!</b>"
    }
}
func Comandos(mensagem string,id_usuario int,username string,ChatID int64,db *sql.DB,bot *tgbotapi.BotAPI) string {
    if strings.Contains(mensagem,"/func_regras"){
        return Func_Regras(mensagem,username,id_usuario)
    }else if strings.Contains(mensagem,"/regras"){
        return Regras()
    }else if strings.Contains(mensagem,"/help"){
        return Help()
    }else if strings.Contains(mensagem,"/admins"){
        return Admins_C(ChatID,bot)
    }else if strings.Contains(mensagem,"/kick"){
        return Kick_Comando(ChatID,mensagem,id_usuario,username,db,bot)
    }else if strings.Contains(mensagem,"/txt_da_vergonha"){
        return Txt_Da_Vergonha()
    }else if strings.Contains(mensagem,"/imagem"){
        return Imagem_D(ChatID,username,mensagem)
    }else if strings.Contains(mensagem,"/clear"){
        return Clear(mensagem,id_usuario,username)
    }else if strings.Contains(mensagem,"/ban"){
        return Ban_Comando(ChatID,mensagem,id_usuario,username,db,bot)
    }
    return "Comando não encontrado,use o /help para saber os comandos!"
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////FUNÇÕES RELACIONADAS AO SISTEMA//////////////////////////////////////
func Leitor(arquivo string) string{
    conteudo_arquivo , err := ioutil.ReadFile(arquivo)
    if err != nil{
        error_check(err)
        return "<b>O arquivo não pode ser lido!Checar logs!</b>"
    }
    return string(conteudo_arquivo)
}
func Escreve(arquivo string,conteudo string){
    os.Remove("imagem.jpg")
    conteudo_arquivo := []byte(conteudo)
    err := ioutil.WriteFile(arquivo, conteudo_arquivo , 0777)
    error_check(err)
}
func Tratamento(mensagem string,comando string)string{
    mensagem = strings.Replace(mensagem,comando,"",-1)
    mensagem = strings.Replace(mensagem," ","",-1)
    if strings.Contains(mensagem,"@"){
        mensagem = strings.Replace(mensagem,"@","",-1)
    }
    return mensagem
}
func jsonS(decode *tgbotapi.User)string{
    a := decode
    out, err := json.Marshal(a)
    error_check(err)
    return string(out)
}
///////////////////////////////////////////"NAVEGADOR"////////////////////////////////////////////////////
func Baixar_Arquivo(url string) string {
  resp, err := http.Get(url)
  error_check(err)
  defer resp.Body.Close()
  b,err := ioutil.ReadAll(resp.Body)
  error_check(err)
  return string(b)
}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func main() {
    bot,err := Inicia_Bot()
    error_check(err)
    db,err := Inicia_Database()
    error_check(err)

    defer db.Close()

    log.Printf("Logado na conta %s", bot.Self.UserName)
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 5 //Tempo pra atualizar

    updates, err := bot.GetUpdatesChan(u)

    for msg := range updates { //Recebe as mensagem
        if msg.Message == nil {
            continue
        }
        //////////////////////////////Verifica o BAN//////////////////
        if msg.Message.NewChatMember != nil{
            jsonBan := jsonS(msg.Message.NewChatMember)
            err := json.Unmarshal([]byte(jsonBan),&banStruct)
            error_check(err)
            if rUser(*db,banStruct.Username,tabela_banidos) == 0{
            }else{
                Kick(banStruct.Id,msg.Message.Chat.ID,bot)
            }
            msg.Message.NewChatMember = nil
        }
        ///////////////////////////////Cadastra o Usuario//////////////////////////////////////
        if rUser(*db,msg.Message.From.UserName,tabela_user) == msg.Message.From.ID{
        }else{
            iUser(*db,msg.Message.From.UserName,msg.Message.From.ID,tabela_user)
            log.Println("O usuario " + msg.Message.From.UserName + "(ID:" + strconv.Itoa(msg.Message.From.ID) + ") foi cadastrado!")
        }
        /////////////////////////////////////////////////////////////////////////////////////////
        log.Printf("[%s] %s", msg.Message.From.UserName, msg.Message.Text)

        if msg.Message.LeftChatMember != nil{
            log.Printf("[%s] foi removido do grupo.",msg.Message.LeftChatMember)
            Mandar_Mensagen(msg.Message.Chat.ID,"<b>Esse deve ter feito merda....( ͡° ͜ʖ ͡°)</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>")
        }
        if msg.Message.NewChatMember != nil{
            log.Printf("[%s] foi adicionado/convidado ao grupo por %s!",msg.Message.NewChatMember,msg.Message.From.UserName)
            Mandar_Mensagen(msg.Message.Chat.ID,"<b>Eai GOleiro , seja bem vindo a alcateia! Mas conta ai , como chegou aqui ?</b>\n\n<b>Versão do Bot:" + Bot_V + "</b>")
        }

        if msg.Message.Text == ""{
            continue
        }else{
            if Verifica_Comando(msg.Message.Text) == true{
                x := Comandos(msg.Message.Text,msg.Message.From.ID,msg.Message.From.UserName,msg.Message.Chat.ID,db,bot)
                Mandar_Mensagen(msg.Message.Chat.ID,x)
            }
        }
    }
}