package main

//	O servidor só precisa gerar o mapa
//	Deve ser implementadas funções que modifiquem o mapa, movem os jogadores, e interajam com os elementos
//
//
//
import (
	"bufio"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"

	"github.com/nsf/termbox-go"
)

var mutex sync.Mutex

func connect() {
	// Conectar ao servidor
	//cria um objeto jogador( ID , ELEMENT )
}

// Define os elementos do jogo
type Elemento struct {
	simbolo  rune
	cor      termbox.Attribute
	corFundo termbox.Attribute
	tangivel bool
}
type Jogador struct {
	ID      int
	Element Elemento
	posX    int
	posY    int
	Online  bool
}
type Servidor struct {
	Jogadores                   [3]Jogador
	mapa                        [][]Elemento
	ultimoElementoSobPersonagem Elemento
	statusMsg                   string
	efeitoNeblina               bool
	revelado                    [][]bool
	raioVisao                   int
	mapa_Inicializado           bool
}

var personagem_1 = Elemento{
	simbolo:  '☺',
	cor:      termbox.ColorRed,
	corFundo: termbox.ColorDefault,
	tangivel: true,
}
var personagem_2 = Elemento{
	simbolo:  '☺',
	cor:      termbox.ColorGreen,
	corFundo: termbox.ColorDefault,
	tangivel: true,
}
var personagem_3 = Elemento{
	simbolo:  '☺',
	cor:      termbox.ColorBlue,
	corFundo: termbox.ColorDefault,
	tangivel: true,
}

// Parede
var parede = Elemento{
	simbolo:  '▤',
	cor:      termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim,
	corFundo: termbox.ColorDarkGray,
	tangivel: true,
}

// Barrreira
var barreira = Elemento{
	simbolo:  '#',
	cor:      termbox.ColorRed,
	corFundo: termbox.ColorDefault,
	tangivel: true,
}

// Vegetação
var vegetacao = Elemento{
	simbolo:  '♣',
	cor:      termbox.ColorGreen,
	corFundo: termbox.ColorDefault,
	tangivel: false,
}

// Elemento vazio
var vazio = Elemento{
	simbolo:  ' ',
	cor:      termbox.ColorDefault,
	corFundo: termbox.ColorDefault,
	tangivel: false,
}

// Elemento para representar áreas não reveladas (efeito de neblina)
var neblina = Elemento{
	simbolo:  '.',
	cor:      termbox.ColorDefault,
	corFundo: termbox.ColorYellow,
	tangivel: false,
}

// Servidor recebe (comando, jogador.posX, jogador.posY, elem, lastElem) e chama a função updatePos

// var mapa [][]Elemento
// var ultimoElementoSobPersonagem = vazio
// var statusMsg string
// var efeitoNeblina = false
// var revelado [][]bool
// var raioVisao int = 3

func (s *Servidor) inicializar() {
	s.ultimoElementoSobPersonagem = vazio
	s.efeitoNeblina = false
	s.raioVisao = 3
	s.mapa_Inicializado = false
	// Inicializa o mapa
	s.carregarMapa("mapa.txt")
	// Inicializa os jogadores
	s.Jogadores[0] = Jogador{ID: 0, Element: personagem_1, posX: -1, posY: -1, Online: false}
	s.Jogadores[1] = Jogador{ID: 1, Element: personagem_2, posX: -1, posY: -1, Online: false}
	s.Jogadores[2] = Jogador{ID: 2, Element: personagem_3, posX: -1, posY: -1, Online: false}
}
func (s *Servidor) sendMapa(trash string, clientMap *[][]Elemento) error { // cliente manda seu mapa, servidor Carrega o mapa
	if s.mapa_Inicializado {
		*clientMap = s.mapa
		return nil
	}
	return fmt.Errorf("Mapa ainda não Inicializado")

}

func (s *Servidor) listenInput(ev string, j *Jogador) { /*TODO*/
	switch ev {
	case "w":
		s.updatePos(j.posX, j.posY-1, j.Element, *j)
	case "a":
		s.updatePos(j.posX-1, j.posY, j.Element, *j)
	case "s":
		s.updatePos(j.posX, j.posY+1, j.Element, *j)
	case "d":
		s.updatePos(j.posX+1, j.posY, j.Element, *j)
	}
}

func (s *Servidor) interact(ev string, j *Jogador) { /*idk what TODO*/

}

// func (s *Servidor) updateMap(trash string, j *Jogador) { /*TODO*/ }
func (s *Servidor) getPlayer(trash string, j *Jogador) error { /*DONE*/
	for i := 0; i < 3; i++ {
		if !s.Jogadores[i].Online {
			s.Jogadores[i].Online = true
			s.Jogadores[i].posX = 10 + i
			s.Jogadores[i].posY = 3

			*j = s.Jogadores[i]
			return nil
		}
	}
	return fmt.Errorf("Não há mais jogadores disponíveis")
}
func main() {
	porta := 8973
	servidor := new(Servidor)
	servidor.inicializar()

	rpc.Register(servidor)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", porta))
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}

	fmt.Println("Servidor aguardando conexões na porta", porta)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}

}

func (s *Servidor) carregarMapa(nomeArquivo string) {
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		panic(err)
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	y := 0
	for scanner.Scan() {
		linhaTexto := scanner.Text()
		var linhaElementos []Elemento
		for x, char := range linhaTexto {
			elementoAtual := vazio
			switch char {
			case parede.simbolo:
				elementoAtual = parede
			case barreira.simbolo:
				elementoAtual = barreira
			case vegetacao.simbolo:
				elementoAtual = vegetacao
			}
			if x == x {
				linhaElementos = append(linhaElementos, elementoAtual)
			}
		}
		s.mapa = append(s.mapa, linhaElementos)
		y++
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	s.mapa_Inicializado = true
}

func (s *Servidor) updatePos(novaPosX int, novaPosY int, elem Elemento, j Jogador) { // Cliente chama essa função para atualizar a posição do elemento
	mutex.Lock()
	if novaPosY >= 0 && novaPosY < len(s.mapa) && novaPosX >= 0 && novaPosX < len(s.mapa[novaPosY]) && s.mapa[novaPosY][novaPosX].tangivel == false {
		s.mapa[j.posY][j.posX] = s.ultimoElementoSobPersonagem     // Restaura o elemento anterior
		s.ultimoElementoSobPersonagem = s.mapa[novaPosY][novaPosX] // Atualiza o elemento sob o personagem
		j.posX, j.posY = novaPosX, novaPosY                        // Move o personagem
		s.mapa[j.posY][j.posX] = j.Element                         // Coloca o personagem na nova posição
	}
	mutex.Unlock()
}