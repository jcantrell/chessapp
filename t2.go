package main

import (
  "errors"
  "net"
  "fmt"
  "bufio"
  "strings"
  "strconv"
//  "encoding/binary"
  "io/ioutil"
  "unicode"
  "os"
)

func checkFullmoves(f string) (int, error) {
  return strconv.Atoi(f)
}

func checkHalfmoves(h string) (int, error) {
  return checkFullmoves(h)
}

func checkFields(s []string) (bool, error) {
  if len(s) != 6 {
    return false, errors.New("Malformed SEN string - too many or too few fields")
  }
  _, errmsg := checkFullmoves(s[5])
  if errmsg != nil {
    return false, errors.New("Malformed fullmoves field")
  }
  _, errmsg = checkHalfmoves(s[4])
  return true, nil
}

func getBitBoards(b string) []byte {
  var bitBoards = make(map[byte]uint64)
  pieceTypes := "rnbqkpRNBQKP"
  var row,col int
  for i := 0; i<len(b); i++ {
    if (b[i] == '/') {
      row++
      col = 0
    } else if (unicode.IsDigit(rune(b[i]))) {
      inc,_ := strconv.Atoi(string(b[i]))
      col += inc
    } else {
      bitBoards[b[i]] |= (1<<(63-(8*row+col)))
      col++
    }
  }

  r := make([]byte,8*12)
  for i := 0; i < len(pieceTypes); i++ {
    for j:=0; j<8; j++ {
      r[i*8+j] = byte(bitBoards[pieceTypes[i]]>>(8*(7-j))) & 0xFF
    }
  }

  return r
}

func main() {
  type ChessState struct {
    boardlayout string
    toplay rune
    castles string
    enpassant string
    halfmove int
    fullmove int
  }

  fmt.Println("Start server...")

  // listen on port 8000
  ln, _ := net.Listen("tcp", ":8000")

  // accept connection
  conn, _ := ln.Accept()

  // run loop forever (or until ctrl-c)
  for {
    // get message, output
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      os.Exit(1)
    }
    words := strings.Fields(message)
    _, errmsg := checkFields(words)
    if errmsg != nil {
      fmt.Println(errmsg)
      continue;
    }
    bf := getBitBoards(words[0])
    ioutil.WriteFile("/tmp/dat1", bf, 0644)
    conn.Write(bf)
    for i := 0; i < len(words); i++ {
      fmt.Println("Message Received:", string(words[i]))
    }
  }
}
