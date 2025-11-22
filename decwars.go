//
// decwars
// This version by Harris Newman
//
// License: This project is licensed under the GNU General Public License v3.0 or later (GPL-3.0-or-later).
// You are free to use, modify, and redistribute this software under the terms of the GPL. 
// There is no warranty of any kind. See the full license text in LICENSE (or COPYING) and 
// the official GPL page: https://www.gnu.org/licenses/gpl-3.0.html.
//
// To do disallow reserved words for username
// reserved words:
//
//
// 1) Install OS
// 2) Install ssh
// 3) Change IP
// 4) Install golang meta package
// 5) apt-get install git
// 6)
// Open your ~/.bashrc file, and enter the following lines at the end:
// export GOPATH=$HOME/gocode:/usr/local/go/bin
// export GOARCH=amd64
// export GOOS=linux
// export GOBIN=$HOME/gocode/bin
// export PATH=$HOME/gocode/bin:$PATH:/usr/local/go/bin
// export GOCERTFILE=~/cert.pem
// export GOKEYFILE=~/key.pem
// export GOHOST=192.168.0.9
//
// 7) go get github.com/mattn/go-sqlite3
// 8)
// download godev
// change debug.go, godev.go, terminal.go and test.go to have golang.org/x/net/websocket instead of code.google.com/p/go.net/websocket
// godev requires certs, howto rebuild (then install cert.pem on chromium):
//  openssl genrsa -des3 -passout pass:x -out server.pass.key 2048
//  openssl rsa -passin pass:x -in server.pass.key -out server.key
//  rm server.pass.key
//  openssl req -new -key server.key -out server.csr
//  openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
//  openssl x509 -in server.crt -out cert.pem
//
// See https://github.com/mattn/go-sqlite3 for DB examples
//
// You need to create a news file wherever you want it
//
// Create a link to the executable: sudo ln /home/newmanh/go/src/decwars/decwars /usr/local/bin/decwars
//
// Reserved word list for usernames:
// * All commands
// * Condition
// * Location
// * Torpedoes
// * Energy
// * Damage
// * Shields
// * Radio
// * Friendly
// * Enemy
//

package main

import (
	"crypto/md5"
	"database/sql"
//	"decwars/telnet"
	"bufio"
	"unicode"
	"bytes"

	"encoding/hex"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	//	"sync"
	"syscall"
	"time"
)

//
// 5 min timeout = 300
//
const timeout = 300 * time.Second

const ctlsp = "\r\f\n "
const ctlonly = "\r\f\n"
const invparm = "Invalid parameter"
const movetbig = "Engineering Officer:  The engines won't take it Captain."
const invstack = "You may not stack the activate command, aborting command line"
const OOR = "Out of range"
const endTract = "Tractor disengaged, captain"
const nAdjPlt = "Not adjacent to planet, Captain"
const nRefuseCapture = "Sir, the enemy refuses our surrender ultimatum!"
const nAPlt = "Not a planet, Captain"
const plCaptured = "Planet captured, Captain"
const nEnRefuse = "Sir, the enemy refuses our help!"
const noTarg = "Phaser control unable to lock on target, Captain."
const highSpeedShieldControl = "High speed shield control activated."
const ABF = "All bases still functional, Captain"
const BPB = " builds planet into base!"
const phasSize = "Weapons Officer:  Improper energy consumption for phaser hit, Captain."
const phasUnable = "Phaser control unable to lock on target, Captain."
const phasFriendL = "Weapons Officer:  Attempting to hit friendly object: "
const phasFriendM = " neutralized by "
const phasOverheating1 = "WARNING! WARNING!  PHASERS OVERHEATING."
const phasOverheating2 = "********** CRACKLE! POP! SIZZLE! POOF! **********"
const phasOverheating3 = "PHASERS DAMAGED."
const phasDamMsg = "Phasers critically damaged."
const torpDamMsg = "Torpedo tubes critically damaged."
const cmpDamMsg = "Computer critically damaged."
const unableRaiseShld = "Unable to raise shields, Captain"
const Destroyed = " is destroyed!"
const TOOR = "Target out of range."
const torpWrongNum = "Insufficient torpedoes for burst!"
const torpHitLong = " makes torpedo hit on "
const bhGulp = " gulps "
const nAdjPort = "Not adjacent to port, Captain"
const endDockMsg = "Undocked, Captain"
const successDock = "Docking completed, Captain"
const endGame = "THE WAR IS OVER!!!"
const endGameWhoC = "The Coalition is victorious!"
const endGameWhoE = "The Empire is victorious!"
const gulpTorp = " gulps torpedo"
const torptoomany = "Weapons Officer:  Captain, we can only launch 3 torpedoes at a time!"
const whatCaptain = "Weapons Officer:  Captain, perhaps you should rest a bit!"

// max len of username
const Usernamemax = 6

// Constants for object side
const SideCoalition = 1
const SideEmpire = 2
const SideNeutral = 4
const SideArcheron = 8

// Constants for side names
const NameCoalition = "Coalition"
const NameEmpire = "Empire"
const NameNeutral = "Neutral"
const NameArcheron = "Archeron"
const NameAll = "All"
const NameFriendlies = "Friendlies"
const NameEnemies = "Enemies"
const NameTargets = "Targets"
const NameStars = "Stars"
const NameShips = "Ships"
const NamePlanets = "Planets"
const NameBases = "Bases"
const NamePorts = "Ports"
const NameClosest = "Closest"
const NameCaptured = "Captured"
const NameEnemy = "Enemy"
const NameFriendly = "Friendly"
const NameSummary = "Summary"
const NameRelative = "Relative"
const NameAbsolute = "Absolute"
const NameBoth = "Both"
const NameComputed = "Computed"
const NameBH = "Black Holes"
const NameLong = "Long"
const NameMedium = "Medium"
const NameShort = "Short"
const NameEnable = "Enable"
const NameDisable = "Disable"
const NameSuon = "su-on"
const NameSuoff = "su-off"
const NameBHoff = "bh-off"
const NameBHon = "bh-on"
const NameARoff = "ar-off"
const NameARon = "ar-on"
const NameBounce = "bounce"
const NameReport = "report"
const NameRemove = "remove"
const NameDelete = "delete"
const NameReset = "reset"
const NameWED = "wed"
const NameIED = "ied"
const NamePTD = "ptd"
const NamePBD = "pbd"
const NameSD = "sd"
const NameCD = "cd"
const NameLSD = "lsd"
const NameRD = "rd"
const NameTBD = "tbd"
const NameTSD = "tsd"
const NamedoHit = "dohit"
const NameTorp = "torpedoes"
const NamePhas = "phasers"
const NameStart = "start"
const NameEndG = "endgame"

//
const NameStarSh = "*"
const NamePlanetSh = "@"
const NameBaseSh = "B"
const NameBHSh = "BH"

//
const NameStarM = "Star"
const NamePlanetM = "Planet"
const NameBaseM = "Base"
const NameBHM = "B-H"

//
const NameStarL = "Star"
const NamePlanetL = "Planet"
const NameBaseL = "Base"
const NameBHL = "Black Hole"

//
const NameOn = "on"
const NameOff = "off"
const NameUp = "up"
const NameDown = "down"
const NameTransfer = "transfer"

// Truths
const On = 1
const Off = 0

// Constants for object type
const TypeShip = 1
const TypePlanet = 2
const TypeStar = 4
const TypeBase = 8
const TypeBH = 16

// Obeject defaults
const InitLifeSup = 5
const InitEnergy = 5000
const InitShield = 2500
const InitPhoTor = 10
const MaxWarpFactor = 6
const MaxImpulseFactor = 1
const MaxBaseRng = 4
const MaxPlanetRng = 2 //like decwar
const MaxArcheronRng = 4
const MaxScanRng = 10

// Constants for IO format
const IOFmtRel = 1
const IOFmtAbs = 2
const IOFmtBoth = 3

// Constants for IO length
const OutLenSh = 1
const OutLenMed = 2
const OutLenLong = 4

// Status constants
const StatGreen = "Green"
const StatYellow = "Yellow"
const StatRed = "Red"
const StatG = 1
const StatY = 2
const StatR = 3

const lenOfString1 = 1
const lenOfString2 = 2
const lenOfString3 = 3
const lenOfString4 = 4

// Constants for delays in milliseconds
const movedelay = 500

// Constants for damage
const WarpEngineDamageFactor = 50
const MaxDam = 300      //max before warp 3 is max
const MaxShipDam = 2500 // if ship damage >= this, you die - decwars=2500
const MaxStarDam = 250  // Max damage for star to explode - was 500
const MaxBaseDam = 500  // Max damage for base to explode - was 1500
const MaxArDam = 200    // Max damage for Archeron - was 400
const KCRIT = 300       //critical damage - all devices
const ImpEngineDamageFactor = 50
const defPhasAmt = 200
const defTorpHitAmt = 300
const defStarAmt = 1000
const defTorpAmt = 3
const defRepRate = 50   // repair rate is 50 not docked, 100 docked
const Basesmax = 9      // FOR TESTING normally 10 (0-9)
const torpFailure = 5   // 10% of the torps misfire
const torpAccuracy = 10 // 10% possibility of inaccuracy per sector shot

// Message types
const msgObjSwallowedBH = 1  //object entered black hole
const msgObjDied = 2         //object died
const msgPhas = 3            //object makes phaser hit
const msgTor = 4             //object makes torp hit
const msgBuildBase = 5       //object makes base
const msgDestroyed = 6       //object is destroyed
const msgStarNova = 7        //star novas
const msgStarDam = 8         //star damages someone
const msgGalaxywideAsst = 9  //galaxy wide call for assistance
const msgGalaxywideBase = 10 //galaxy wide notify base destroyed
const msgRomDetected = 11    //Rom in area
const msgShp2Shp = 12        //energy transfer
const msgTractorAct = 13     //TractorBeam activated
const msgTractorBroken = 14  //TractorBeam broken
const msgTorpNeu = 15        //Torp neutralized
const msgTorpMisfire = 16    //torpedo misfires
const msgTorpHit = 17        //torpedo hit
const msgTorpMiss = 18       //torpedo miss
const msgStar = 19           //Star hit message
const msgGulp = 20           //black hole ate something
const msgDisp = 21           //object displaced
const msgEndGame = 22        //announce end game
const msgEndCoalition = 23   //coalition won
const msgEndEmpire = 24      //empire won
const msgGulpTorp = 25       //black hole ate torpedo

// Device bits for hit calc
const Warp = 1
const ImpEng = 2
const PhoTor = 4
const Phas = 8
const ShldS = 16
const Cmp = 32
const LifeSup = 64
const Radio = 128
const Tractor = 256

// Type of hit
const phasHit = 1
const torpHit = 2

// Chance of archeron doing a move
const archeronPercent = 15

// Chance of planet loosing a build
const buildPercent = 50

// Mutex declaration
// var mu = &sync.Mutex{}

//
// Sqlite3 objectdb mutex
//
//var sqlite3mu = &sync.Mutex{}

//
// Startup time
//
var startupTime = time.Now()

//
// Game time - initial setting
//
var gameTime = time.Now()

//
// Game number - initial setting
//
var GameNumber = 1

type Loc struct {
	Vpos string
	Hpos string
}

//
// Username/address array
//
type Constr struct {
	Username      string
	Connection    net.Conn
	Remoteaddress string
}

// map the name to Constr structure
var Conmap map[string]Constr

//
// Tuneable variables for the universe, with defaults
//
// max size of universe
//
var Vmax = 74
var Hmax = 74
var Starsmax = 120
var BHsmax = 120
var Planetsmax = 60
var Archeronmax = 5
var Portnum = 1701
//
// telnet stuff
//
const (
	CR = byte('\r')
	LF = byte('\n')
)

const (
	cmdSE   = 240
	cmdNOP  = 241
	cmdData = 242

	cmdBreak = 243
	//by hsn
	cmdABT  = 238
	cmdEOR  = 239
	cmdAYT  = 246
	cmdCTLC = 244
	cmdAO   = 245

	// end hsn
	cmdGA = 249
	cmdSB = 250

	cmdWill = 251
	cmdWont = 252
	cmdDo   = 253
	cmdDont = 254

	cmdIAC = 255
)

const (
	optEcho            = 1
	optSuppressGoAhead = 3
	//	optTerminalType    = 24
	optNAWS = 31
)

// Conn implements net.Conn interface for Telnet protocol plus some set of
// Telnet specific methods.
type Conn struct {
	net.Conn
	r *bufio.Reader

	unixWriteMode bool

	cliSuppressGoAhead bool
	cliEcho            bool
}

func NewConn(conn net.Conn) (*Conn, error) {
	c := Conn{
		Conn: conn,
		r:    bufio.NewReaderSize(conn, 256),
	}
	return &c, nil
}

func Dial(network, addr string) (*Conn, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

func DialTimeout(network, addr string, timeout time.Duration) (*Conn, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return NewConn(conn)
}

// SetUnixWriteMode sets flag that applies only to the Write method.
// If set, Write converts any '\n' (LF) to '\r\n' (CR LF).
func (c *Conn) SetUnixWriteMode(uwm bool) {
	c.unixWriteMode = uwm
}

func (c *Conn) do(option byte) error {
	//log.Println("do:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdDo, option})
	return err
}

func (c *Conn) dont(option byte) error {
	//log.Println("dont:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdDont, option})
	return err
}

func (c *Conn) will(option byte) error {
	//log.Println("will:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdWill, option})
	return err
}

func (c *Conn) wont(option byte) error {
	//log.Println("wont:", option)
	_, err := c.Conn.Write([]byte{cmdIAC, cmdWont, option})
	return err
}

func (c *Conn) sub(opt byte, data ...byte) error {
	if _, err := c.Conn.Write([]byte{cmdIAC, cmdSB, opt}); err != nil {
		return err
	}
	if _, err := c.Write(data); err != nil {
		return err
	}
	_, err := c.Conn.Write([]byte{cmdIAC, cmdSE})
	return err
}

func (c *Conn) deny(cmd, opt byte) (err error) {
	switch cmd {
	case cmdDo:
		err = c.wont(opt)
	case cmdDont:
		// nop
	case cmdWill, cmdWont:
		err = c.dont(opt)
	}
	return
}

func (c *Conn) skipSubneg() error {
	for {
		if b, err := c.r.ReadByte(); err != nil {
			return err
		} else if b == cmdIAC {
			if b, err = c.r.ReadByte(); err != nil {
				return err
			} else if b == cmdSE {
				return nil
			}
		}
	}
}

func (c *Conn) cmd(cmd byte) error {
	switch cmd {
	case cmdGA:
		return nil
		//
		//  HSN - test for ayt
		//
	case cmdABT:
		return nil

	case cmdAYT:
		return nil

	case cmdCTLC:
		return nil

	case cmdEOR:
		return nil

	case cmdAO:
		return nil

	case cmdNOP:
		return nil
	//
	// end HSN changes
	//
	case cmdDo, cmdDont, cmdWill, cmdWont:
		// Process cmd after this switch.
	case cmdSB:
		return c.skipSubneg()
	default:
		return fmt.Errorf("unknown command: %d", cmd)
	}
	// Read an option
	o, err := c.r.ReadByte()
	if err != nil {
		return err
	}
// fmt.Println("received cmd:", cmd, "option:", o)
	switch o {
	case optEcho:
		// Accept any echo configuration.
		switch cmd {
		case cmdDo:
//			if !c.cliEcho {
				c.cliEcho = false
//				err = c.will(o)
// hsn fix putty issue:
				err = c.wont(o)
//
//			}
		case cmdDont:
//			if c.cliEcho {
				c.cliEcho = false
				err = c.wont(o)
//			}
		case cmdWill:
			err = c.do(o)
// hsn fix putty issue:
//			err = c.dont(o)
//
		case cmdWont:
//			err = c.dont(o)
			err = c.do(o)
		}
	case optSuppressGoAhead:
		// We don't use GA so can allways accept every configuration
		switch cmd {
		case cmdDo:
			if !c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = true
				err = c.will(o)
			}
		case cmdDont:
			if c.cliSuppressGoAhead {
				c.cliSuppressGoAhead = false
				err = c.wont(o)
			}
		case cmdWill:
			err = c.do(o)
		case cmdWont:
			err = c.dont(o)

		}
	case optNAWS:
		if cmd != cmdDo {
			err = c.deny(cmd, o)
			break
		}
		if err = c.will(o); err != nil {
			break
		}
		// Reply with max window size: 65535x65535
		err = c.sub(o, 255, 255, 255, 255)
	default:
		// Deny any other option
		err = c.deny(cmd, o)
	}
	return err
}

func (c *Conn) tryReadByte() (b byte, retry bool, err error) {
	b, err = c.r.ReadByte()
	if err != nil || b != cmdIAC {
		return
	}
	b, err = c.r.ReadByte()
	if err != nil {
		return
	}
	if b != cmdIAC {
		err = c.cmd(b)
		if err != nil {
			return
		}
		retry = true
	}
	return
}

// SetEcho tries to enable/disable echo on server side. Typically telnet
// servers doesn't support this.
func (c *Conn) SetEcho(echo bool) error {
	if echo {
		return c.do(optEcho)
	}
	return c.dont(optEcho)
}

// ReadByte works like bufio.ReadByte
func (c *Conn) ReadByte() (b byte, err error) {
	retry := true
	for retry && err == nil {
		b, retry, err = c.tryReadByte()
	}
	return
}

// ReadRune works like bufio.ReadRune
func (c *Conn) ReadRune() (r rune, size int, err error) {
loop:
	r, size, err = c.r.ReadRune()
	if err != nil {
		return
	}
	if r != unicode.ReplacementChar || size != 1 {
		// Properly readed rune
		return
	}
	// Bad rune
	err = c.r.UnreadRune()
	if err != nil {
		return
	}
	// Read telnet command or escaped IAC
	_, retry, err := c.tryReadByte()
	if err != nil {
		return
	}
	if retry {
		// This bad rune was a begining of telnet command. Try read next rune.
		goto loop
	}
	// Return escaped IAC as unicode.ReplacementChar
	return
}

// Read is for implement an io.Reader interface
func (c *Conn) Read(buf []byte) (int, error) {
	var n int
	for n < len(buf) {
		b, err := c.ReadByte()
		if err != nil {
			return n, err
		}
		//log.Printf("char: %d %q", b, b)
		buf[n] = b
		n++
		if c.r.Buffered() == 0 {
			// Try don't block if can return some data
			break
		}
	}
	return n, nil
}

// ReadBytes works like bufio.ReadBytes
func (c *Conn) ReadBytes(delim byte) ([]byte, error) {
	var line []byte
	for {
		b, err := c.ReadByte()
		if err != nil {
			return nil, err
		}

		line = append(line, b)
		if b == delim {
			break
		}
	}
	return line, nil
}

// SkipBytes works like ReadBytes but skips all read data.
func (c *Conn) SkipBytes(delim byte) error {
	for {
		b, err := c.ReadByte()
		if err != nil {
			return err
		}
		if b == delim {
			break
		}
	}
	return nil
}

// ReadString works like bufio.ReadString
func (c *Conn) ReadString(delim byte) (string, error) {
	bytes, err := c.ReadBytes(delim)
	return string(bytes), err
}

func (c *Conn) readUntil(read bool, delims ...string) ([]byte, int, error) {
	if len(delims) == 0 {
		return nil, 0, nil
	}
	p := make([]string, len(delims))
	for i, s := range delims {
		if len(s) == 0 {
			return nil, 0, nil
		}
		p[i] = s
	}
	var line []byte
	for {
		b, err := c.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		if read {
			line = append(line, b)
		}
		for i, s := range p {
			if s[0] == b {
				if len(s) == 1 {
					return line, i, nil
				}
				p[i] = s[1:]
			} else {
				p[i] = delims[i]
			}
		}
	}
	panic(nil)
}

// ReadUntilIndex reads from connection until one of delimiters occurs. Returns
// read data and an index of delimiter or error.
func (c *Conn) ReadUntilIndex(delims ...string) ([]byte, int, error) {
	return c.readUntil(true, delims...)
}

// ReadUntil works like ReadUntilIndex but don't return a delimiter index.
func (c *Conn) ReadUntil(delims ...string) ([]byte, error) {
	d, _, err := c.readUntil(true, delims...)
	return d, err
}

// SkipUntilIndex works like ReadUntilIndex but skips all read data.
func (c *Conn) SkipUntilIndex(delims ...string) (int, error) {
	_, i, err := c.readUntil(false, delims...)
	return i, err
}

// SkipUntil works like ReadUntil but skips all read data.
func (c *Conn) SkipUntil(delims ...string) error {
	_, _, err := c.readUntil(false, delims...)
	return err
}

// Write is for implement an io.Writer interface
func (c *Conn) Write(buf []byte) (int, error) {
	search := "\xff"
	if c.unixWriteMode {
		search = "\xff\n"
	}
	var (
		n   int
		err error
	)
	for len(buf) > 0 {
		var k int
		i := bytes.IndexAny(buf, search)
		if i == -1 {
			k, err = c.Conn.Write(buf)
			n += k
			break
		}
		k, err = c.Conn.Write(buf[:i])
		n += k
		if err != nil {
			break
		}
		switch buf[i] {
		case LF:
			k, err = c.Conn.Write([]byte{CR, LF})
		case cmdIAC:
			k, err = c.Conn.Write([]byte{cmdIAC, cmdIAC})
		}
		n += k
		if err != nil {
			break
		}
		buf = buf[i+1:]
	}
	return n, err
}
//
//end telnet stuff
//
// Makes atomic transactions
//
func atomicBegin(objectsdb *sql.DB) (*sql.Tx, error) {
	var err error
	//	sqlite3// mu.Lock()
	tx, err := objectsdb.Begin()
	return tx, err
}

func atomicCommit(tx *sql.Tx) error {
	var err error
	err = tx.Commit()
	//	sqlite3// mu.Unlock()
	return err
}

//
// Replenish energy based on type of port you are docked to
//  cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc
// Resource                        Base    Planet
// ----------------------------------------------
// Ship energy                    +1000      +500
// Shield energy                   +500      +250
// Photon Torpedoes                 +10        +5
// Life Support Reserves             +5        +5
// Ship Damage                     -100       -50
// Ship Damage, if already docked  -200      -100
//
func doReplenishment(username string, myDockFlag string, objectsdb *sql.DB) {
	var portObjtype int
	var myShipEnergy int
	var myShld int
	var myPhoTor int
	var myLifeSup int
	var myShipDam int
	objectsdb.QueryRow("select Objtype from objects WHERE Nme = ?", myDockFlag).Scan(&portObjtype)
	objectsdb.QueryRow("select ShipEnergy, Shld, PhoTor, LifeSup, ShipDam from objects WHERE Nme = ?", username).Scan(&myShipEnergy, &myShld, &myPhoTor, &myLifeSup, &myShipDam)
	if portObjtype == TypePlanet {
		myShipEnergy = myShipEnergy + 500
		myShld = myShld + 250
		myPhoTor = myPhoTor + 5
		myLifeSup = myLifeSup + 5
		myShipDam = myShipDam - 50
	} else {
		myShipEnergy = myShipEnergy + 1000
		myShld = myShld + 500
		myPhoTor = myPhoTor + 10
		myLifeSup = myLifeSup + 5
		myShipDam = myShipDam - 100
	}
	// Can only be to max levels...
	if myShipEnergy > InitEnergy {
		myShipEnergy = InitEnergy
	}
	if myShld > InitShield {
		myShld = InitShield
	}
	if myPhoTor > InitPhoTor {
		myPhoTor = InitPhoTor
	}
	if myLifeSup > InitLifeSup {
		myLifeSup = InitLifeSup
	}
	if myShipDam < 0 {
		myShipDam = 0
	}
	// Update the database for the new values
	objectsdb.Exec("UPDATE objects set ShipEnergy = ?, Shld = ?, PhoTor = ?, LifeSup = ?, ShipDam = ? where Nme = ?", &myShipEnergy, &myShld, &myPhoTor, &myLifeSup, &myShipDam, username)
}

//
// do a delay to penalize user for doing a command that takes time
//
func doGameDelay(username string, nummoves int, objectsdb *sql.DB) {
	var myDockFlag string
	time.Sleep(time.Duration(nummoves) * time.Duration(movedelay) * time.Millisecond)

	//
	// Do auto repairs - which is 30 and done on every move
	//
	doRepair(username, 30, objectsdb)

	//
	// Do replenishment - only if docked
	//
	objectsdb.QueryRow("select DockFlag from objects WHERE Nme = ?", username).Scan(&myDockFlag)
	if myDockFlag != "" {
		doReplenishment(username, myDockFlag, objectsdb)
	}
}

//
// count the number of bits in an int
//
func popcount(input uint64) uint64 {

	var mask1, mask2, mask3, mask4, mask5, mask6 uint64
	mask1 = 6148914691236517205 // 01010101...
	mask2 = 3689348814741910323 // 00110011...
	mask3 = 1085102592571150095 // 00001111...
	mask4 = 71777214294589695   // 8 zeroes, 8 ones, etc...
	mask5 = 70367670468607      // 16 zeroes, 16 ones, etc...
	mask6 = 4294967295          // 32 zeroes, 32 ones

	input = (input & mask1) + ((input >> 1) & mask1)
	input = (input & mask2) + ((input >> 2) & mask2)
	input = (input & mask3) + ((input >> 4) & mask3)
	input = (input & mask4) + ((input >> 8) & mask4)
	input = (input & mask5) + ((input >> 16) & mask5)
	input = (input & mask6) + ((input >> 32) & mask6)

	return uint64(input)
}

//
// Displace an object away from the hitterLocx/y & notify
//
func displacement(whoNme string, objectsdb *sql.DB, conn *Conn, hitterLocx int, hitterLocy int, pointsdb *sql.DB, usersdb *sql.DB) {
	var whoLocx int
	var whoLocy int
	var whoType int
	var dispNme string
	var dispType int
	var newLocx int
	var newLocy int
	var err error
	var locxdisplacement int
	var locydisplacement int
	var myTractorOn int
	var myDockFlag string

	objectsdb.QueryRow("select Locx, Locy, Objtype from objects WHERE Nme = ?", whoNme).Scan(&whoLocx, &whoLocy, &whoType)
	//
	// if the type is bh or star there can be no displacement
	//
	if whoType == TypeBH || whoType == TypeStar {
		return
	}

	//
	// Compute displacement
	//
	if (hitterLocx - whoLocx) != 0 {
		locxdisplacement = (hitterLocx - whoLocx) / Abs(hitterLocx-whoLocx) * -1
	} else {
		locxdisplacement = 0
	}

	if (hitterLocy - whoLocy) != 0 {
		locydisplacement = (hitterLocy - whoLocy) / Abs(hitterLocy-whoLocy) * -1
	} else {
		locydisplacement = 0
	}

	newLocx = whoLocx + locxdisplacement
	newLocy = whoLocy + locydisplacement

	if newLocx < 0 {
		newLocx = 0
	}
	if newLocy < 0 {
		newLocy = 0
	}
	if newLocx > Vmax {
		newLocx = Vmax
	}
	if newLocy > Hmax {
		newLocy = Hmax
	}
	//
	// Is there something where you are displaced to? if so, if it's a black hole you die, otherwise return
	//
	err = objectsdb.QueryRow("select Nme, ObjType from objects WHERE Locx = ? and Locy = ?", newLocx, newLocy).Scan(&dispNme, &dispType)
	if err == nil { //something in the way
		if dispType == TypeBH { //black hole in the way

			//
			//	Are we being tractored or are we docked - if so break the relationship before the move
			//
			objectsdb.QueryRow("select TractorOn, DockFlag from objects WHERE Nme = ?", whoNme).Scan(&myTractorOn, &myDockFlag)
			if myTractorOn == 1 {
				endTractorOn(whoNme, conn, objectsdb)
			}
			if myDockFlag != "" {
				endDock(whoNme, conn, objectsdb)
			}
			notify(whoNme, conn, newLocx, newLocy, msgDisp, dispNme, newLocx, newLocy, objectsdb, 0, 0)

			notify(whoNme, conn, whoLocx, whoLocy, msgGulp, dispNme, newLocx, newLocy, objectsdb, 0, 0)
			dbDelobjects(conn, objectsdb, whoNme, pointsdb, usersdb)
			return
		}
		return
	} else { // nothing in the way
		//
		//	Are we being tractored - if so break the tractor beam before the move
		//
		objectsdb.QueryRow("select TractorOn, DockFlag from objects WHERE Nme = ?", whoNme).Scan(&myTractorOn, &myDockFlag)
		if myTractorOn == 1 {
			endTractorOn(whoNme, conn, objectsdb)
		}
		if myDockFlag != "" {
			endDock(whoNme, conn, objectsdb)
		}
		objectsdb.Exec("UPDATE objects set Locx = ?, Locy = ? where Nme = ?", newLocx, newLocy, whoNme)
		notify(whoNme, conn, whoLocx, whoLocy, msgDisp, whoNme, newLocx, newLocy, objectsdb, 0, 0)
		return
	}
}

//
// Function to update the points db based on incident
//
func doPointsUpdate(hitter string, hitterSide int, whoNme string, whoSide, objtype int, sharedPhasSize int, pointsdb *sql.DB, usersdb *sql.DB) {
	//
	// Each update requires updates to both pointsdb and userdb
	//
	// check for damage to bases
	//
	if objtype == TypeBase {
		pointsdb.Exec("UPDATE points set DamToBases = DamToBases+? where Nme = ?", sharedPhasSize, hitter)
		usersdb.Exec("UPDATE users set DamToBases = DamToBases+? where name=?", sharedPhasSize, hitter)
	}
	//
	// check for damage to ships
	//
	if objtype == TypeShip {
		pointsdb.Exec("UPDATE points set DamToShip = DamToShip+? where Nme = ?", sharedPhasSize, hitter)
		usersdb.Exec("UPDATE users set DamToShip = DamToShip+? where name=?", sharedPhasSize, hitter)
	}
	//
	// check for damage to stars
	//
	if objtype == TypeStar {
		pointsdb.Exec("UPDATE points set DamToStars = DamToStars+? where Nme = ?", sharedPhasSize, hitter)
		usersdb.Exec("UPDATE users set DamToStars = DamToStars+? where name=?", sharedPhasSize, hitter)
	}
	//
	// check for damage to planets
	//
	if objtype == TypePlanet {
		pointsdb.Exec("UPDATE points set DamToPlanets = DamToPlanets+? where Nme = ?", sharedPhasSize, hitter)
		usersdb.Exec("UPDATE users set DamToPlanets = DamToPlanets+? where name=?", sharedPhasSize, hitter)
	}
	//
	// check for damage to bh
	//
	if objtype == TypeBH {
		pointsdb.Exec("UPDATE points set DamToBH = DamToBH+? where Nme = ?", sharedPhasSize, hitter)
		usersdb.Exec("UPDATE users set DamToBH = DamToBH+? where name=?", sharedPhasSize, hitter)
	}
}

//
// Inflict damages based on a "hit" - randomly hit devices based on an impact amount - notify if "knocked down"
// 1) There are 9 devices that can be hit. Create a random uint16 that will identify which devices were hit
// 2) Count the bits that are set on the random uint16, that is the # of devices hit.  Split the hit equally
// 3) For each device hit, hit x% of the equal hit
// 4) The last device gets x% + the remainder
// 5) hitterlocx/y is used for star explosions and determining displacement
// 6) Update the points accordingly
//
func doHit(whoNme string, whoSide int, Locx int, Locy int, Shld int, ShldUp int, WarpEngDam int, ImpEngDam int, PhoTorDam int, PhasDam int, ShldDam int, CmpDam int, LifeSupDam int, RadioDam int, TractorDam int, ShipDam int, Objtype int, hitSize int, objectsdb *sql.DB, conn *Conn, hitType int, hitter string, hitterLocx int, hitterLocy int, hitterSide int, pointsdb *sql.DB, usersdb *sql.DB) {
	var whoLocx int
	var whoLocy int
	var sharedPhasSize int
	var whichDevHit uint64
	var numBits uint64
	var starhitLocx int
	var starhitLocy int
	var starhitNme string
	var starhitSide int
	var err error
	var starhitShld int
	var starhitShldUp int
	var starhitWarpEngDam int
	var starhitImpEngDam int
	var starhitPhoTorDam int
	var starhitPhasDam int
	var starhitShldDam int
	var starhitCmpDam int
	var starhitLifeSupDam int
	var starhitRadioDam int
	var starhitTractorDam int
	var starhitShipDam int
	var starhitObjtype int
	//
	// Points stuff for pointsdb
	//
	var myDamToBH int
	var myDamToBases int
	var myDamToShip int
	var myDamToStars int
	var myDamToPlanets int
	var myNumOfShips int
	var myNumOfStarDates int
	//
	// Points stuff for usersdb
	//
	var uDamToBases int
	//	var uDamToBH int
	var uDamToShip int
	var uDamToStars int
	var uDamToPlanets int
	var uNumOfShips int
	var uNumOfStarDates int
	var bld int

	// fmt.Println("* dohit called with: whoNme:", whoNme, " whoSide:", whoSide, " Locx:", Locx , " Locy:", Locy , " Shld:", Shld , " ShldUp:", ShldUp , " WarpEngDam:", WarpEngDam , " ImpEngDam:", ImpEngDam , " PhoTorDam:", PhoTorDam , " PhasDam:", PhasDam , " ShldDam:", ShldDam , " CmpDam:", CmpDam , " LifeSupDam:", LifeSupDam , " RadioDam:", RadioDam , " TractorDam:", TractorDam , " ShipDam:", ShipDam , " Objtype:", Objtype , " hitSize:", hitSize , " Hittype:", hitType , " Hitter:", hitter , " Hitterlocx:", hitterLocx , " HitterLocy:", hitterLocy , " hitterside:", hitterSide )

	if whoNme == hitter {
		return
	}
	//
	// Step 1 - get the data to update for the user's points
	//
	err = pointsdb.QueryRow("select DamToBH, DamToBases, DamToShip, DamToStars, DamToPlanets, NumOfShips, NumOfStarDates  from points where Nme=?", hitter).Scan(&myDamToBH, &myDamToBases, &myDamToShip, &myDamToStars, &myDamToPlanets, &myNumOfShips, &myNumOfStarDates)

	err = usersdb.QueryRow("select DamToShip, DamToStars, DamToPlanets, NumOfShips, NumOfStarDates from users where name=?", hitter).Scan(&uDamToBases, &uDamToShip, &uDamToStars, &uDamToPlanets, &uNumOfShips, &uNumOfStarDates)

	//
	// Step 2:
	// if shields are 100% and up, they get the full hit
	//
	/*	if ShldUp == On {
				if Shld == InitShield {
					Shld = Shld - hitSize
					ShldDam = ShldDam + hitSize
					if ShldDam >= MaxDam {
						ShldUp = Off
					}
					objectsdb.Exec("UPDATE objects set Shld = ?, ShldDam = ?, ShipDam = ?, ShldUp = ?  WHERE Nme = ?", Shld, ShldDam, ShipDam, ShldUp, whoNme)
					if hitType == torpHit {
						displacement(whoNme, objectsdb, conn, hitterLocx, hitterLocy, pointsdb, usersdb)
					}
		//			return
				}
			}
	*/
	//
	// Randomly impact each device based on impactAmt & if your shields are up or down
	// Shields not full or down, all devices can get hit (shields not full).
	// Calculate # devices hit, and which
	//
	//	whichDevHit = uint64(rand.Intn(512))
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	a := r1.Intn(512)
	whichDevHit = uint64(a)

	//
	// How many bits in whichDevHit are set? If zero make it all
	//

	numBits = popcount(whichDevHit)

	if numBits == 0 {
		numBits = 512
	}

	sharedPhasSize = int(hitSize) / int(numBits)

	// fmt.Println("** sharedPhasSize:", sharedPhasSize, " hitSize:",hitSize, " numbits:",numBits, " ShldUp:", ShldUp, " ShldDam:", ShldDam, "whiDevHit & shlds:", (whichDevHit & ShldS))
	// if Shlds up and full they only take the hit
	// if shlds up and partial, based on rand # * shield percent
	if ShldUp == On {
		if Shld == InitShield {
			Shld = Shld - sharedPhasSize
			ShldDam = ShldDam + sharedPhasSize
			if ShldDam >= MaxDam {
				ShldUp = Off
			}
		} else {
			if (whichDevHit & ShldS) > 0 {
				ShldDam1 := int(float64(sharedPhasSize) * calcShields(Shld) / 100)
				ShldDam = ShldDam + ShldDam1
				Shld = Shld - ShldDam1
				// fmt.Println("*****shields up, not max, shlddam:",ShldDam, " shld:",Shld)
				if ShldDam >= MaxDam {
					ShldUp = Off
				}
			}
			if (whichDevHit & Warp) > 0 {
				WarpEngDam = WarpEngDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
			if (whichDevHit & ImpEng) > 0 {
				ImpEngDam = ImpEngDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
			if (whichDevHit & PhoTor) > 0 {
				PhoTorDam = PhoTorDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
			if (whichDevHit & Phas) > 0 {
				PhasDam = PhasDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
			if (whichDevHit & Cmp) > 0 {
				CmpDam = CmpDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
			if (whichDevHit & LifeSup) > 0 {
				LifeSupDam = LifeSupDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
			if (whichDevHit & Radio) > 0 {
				RadioDam = RadioDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
			if (whichDevHit & Tractor) > 0 {
				TractorDam = TractorDam + int(float64(sharedPhasSize)*calcShields(Shld)/100)
			}
		}
	} else {
		//
		// If shields are down, get the full hit
		//
		if (whichDevHit & Warp) > 0 {
			WarpEngDam = WarpEngDam + sharedPhasSize
		}
		if (whichDevHit & ImpEng) > 0 {
			ImpEngDam = ImpEngDam + sharedPhasSize
		}
		if (whichDevHit & PhoTor) > 0 {
			PhoTorDam = PhoTorDam + sharedPhasSize
		}
		if (whichDevHit & Phas) > 0 {
			PhasDam = PhasDam + sharedPhasSize
		}
		if (whichDevHit & ShldS) > 0 {
			ShldDam = ShldDam + sharedPhasSize
			Shld = Shld - sharedPhasSize
			// fmt.Println("*****shields down, not max, shlddam:",ShldDam, " Shld:", Shld)
		}
		if (whichDevHit & Cmp) > 0 {
			CmpDam = CmpDam + sharedPhasSize
		}
		if (whichDevHit & LifeSup) > 0 {
			LifeSupDam = LifeSupDam + sharedPhasSize
		}
		if (whichDevHit & Radio) > 0 {
			RadioDam = RadioDam + sharedPhasSize
		}
		if (whichDevHit & Tractor) > 0 {
			TractorDam = TractorDam + sharedPhasSize
		}
	}
	// ship damage is increased by hitsize
	ShipDam = ShipDam + hitSize

	// fmt.Println("*** update:", Shld, ShldDam, ShipDam, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, CmpDam, LifeSupDam, TractorDam, ShldUp, RadioDam, whoNme)

	_, err = objectsdb.Exec("UPDATE objects set Shld = ?, ShldDam = ?, ShipDam = ?, WarpEngDam = ?, ImpEngDam = ?, PhoTorDam = ?, PhasDam = ?, CmpDam = ?, LifeSupDam = ?, TractorDam = ?, ShldUp = ?, RadioDam = ?  WHERE Nme = ?",
		Shld, ShldDam, ShipDam, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, CmpDam, LifeSupDam, TractorDam, ShldUp, RadioDam, whoNme)

	// fmt.Println("###### error on update:",err)
	//
	// Tell object that they got damage to their devices if they became critical
	//
	if Conmap[whoNme].Connection != nil {
		if Objtype == TypeShip {
			if WarpEngDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our warp engines are critically damaged!")
			}
			if ImpEngDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our impulse engines are critically damaged!")
			}
			if PhoTorDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our photon torpedos are critically damaged!")
			}
			if PhasDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our phaser banks are critically damaged!")
			}
			if ShldDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our shields are critically damaged!")
			}
			if CmpDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our computer is critically damaged!")
			}
			if LifeSupDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our life support is critically damaged!")
			}
			if RadioDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our radio is critically damaged!")
			}
			if TractorDam >= MaxDam {
				sendln(Conmap[whoNme].Connection.(*Conn), "Captain, our tractor beam is critically damaged!")
			}
		}
	}
	//
	// If the hit type is a torpedo, the star will explode if it' damage > max allowed
	// Exploding stars can hit other objects nearby - hit is same as a torpedo
	// need to add displacement randomly happening for star hits
	if hitType == torpHit {
		//
		// Notify everyone of the hit 06/10/20
		//
		notify(hitter, conn, hitterLocx, hitterLocy, msgTorpHit, whoNme, whoLocx, whoLocy, objectsdb, hitSize, 0)
		//

		displacement(whoNme, objectsdb, conn, hitterLocx, hitterLocy, pointsdb, usersdb)
		//test of killing planets or bases
		if Objtype == TypeStar {
			if ShipDam >= MaxStarDam { // star explodes
				//orig				notify(whoNme, conn, whoLocx, whoLocy, msgStarNova, whoNme, whoLocx, whoLocy, objectsdb, 0, 0)
				notify(whoNme, conn, Locx, Locy, msgStarNova, whoNme, Locx, Locy, objectsdb, 0, 0)
				// fmt.Println("got into star nova section, star dam:",ShipDam ," max dam:" ,MaxStarDam,"objtype:", Objtype, "notify parms: ",whoNme, whoLocx, whoLocy, msgStarNova, whoNme, whoLocx, whoLocy, 0, 0,Locx, Locy)
				dbDelobjects(conn, objectsdb, whoNme, pointsdb, usersdb)
				//
				// if star is adjacent to object, hit the adjacent object/hit is same as torpedo
				// sqlite sucks so have to do the search around the star one at a time to recurse
				//
				err = objectsdb.QueryRow("select Nme, Side, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(Locx+1)+" and Locy ="+strconv.Itoa(Locy-1)).Scan(&starhitNme, &starhitSide, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err1:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}
				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(Locx+1)+" and Locy ="+strconv.Itoa(Locy)).Scan(&starhitNme, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err2:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}

				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(Locx+1)+" and Locy ="+strconv.Itoa(Locy+1)).Scan(&starhitNme, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err3:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}
				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(Locx)+" and Locy ="+strconv.Itoa(Locy-1)).Scan(&starhitNme, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err4:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}

				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(Locx)+" and Locy ="+strconv.Itoa(Locy+1)).Scan(&starhitNme, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err5:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}
				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(Locx-1)+" and Locy ="+strconv.Itoa(Locy-1)).Scan(&starhitNme, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err6:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}

				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(whoLocx-1)+" and Locy ="+strconv.Itoa(whoLocy)).Scan(&starhitNme, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err7:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}
				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype from objects where Locx="+strconv.Itoa(whoLocx-1)+" and Locy ="+strconv.Itoa(whoLocy+1)).Scan(&starhitNme, &starhitLocx, &starhitLocy, &starhitShld, &starhitShldUp, &starhitWarpEngDam, &starhitImpEngDam, &starhitPhoTorDam, &starhitPhasDam, &starhitShldDam, &starhitCmpDam, &starhitLifeSupDam, &starhitRadioDam, &starhitTractorDam, &starhitShipDam, &starhitObjtype)
				//				// fmt.Println("err8:", err, Locx, Locy, starhitNme, starhitLocx, starhitLocy)
				if err == nil {
					notify(whoNme, conn, whoLocx, whoLocy, msgStar, starhitNme, starhitLocx, starhitLocy, objectsdb, defStarAmt, 0)
					doHit(starhitNme, starhitSide, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, defStarAmt, objectsdb, conn, torpHit, whoNme, whoLocx, whoLocy, whoSide, pointsdb, usersdb)
				}
			}
		} else {
			if Objtype == TypeBase {
				if ShipDam >= MaxBaseDam {
					notify(whoNme, conn, whoLocx, whoLocy, msgStarNova, whoNme, whoLocx, whoLocy, objectsdb, 0, 0)
					dbDelobjects(conn, objectsdb, whoNme, pointsdb, usersdb)
				}
			}
			if Objtype == TypeBH {
				notify(whoNme, conn, whoLocx, whoLocy, msgGulpTorp, whoNme, whoLocx, whoLocy, objectsdb, 0, 0)
			}
			if Objtype == TypePlanet {
				if ShipDam >= MaxBaseDam {
					notify(whoNme, conn, whoLocx, whoLocy, msgStarNova, whoNme, whoLocx, whoLocy, objectsdb, 0, 0)
					dbDelobjects(conn, objectsdb, whoNme, pointsdb, usersdb)
				}
				//
				// Torpedos can always reduce the # builds on a planet
				//
				objectsdb.QueryRow("select Builds from objects WHERE Nme = ?", whoNme).Scan(&bld)
				if bld > 0 {
					objectsdb.Exec("Update objects set Builds = Builds - 1 where Nme = ?", whoNme)
				}
			}
		}
		//
		// Notify everyone of the hit of a phaser and reduce the builds if necessary
		//
	} else {
		notify(hitter, conn, hitterLocx, hitterLocy, msgPhas, whoNme, whoLocx, whoLocy, objectsdb, hitSize, 0)
		//
		// If planet has builds reduce by 1 randomly (buildPercent chance)
		//
		r1 := rand.New(s1)
		a := r1.Intn(100)
		if a <= buildPercent && Objtype == TypePlanet {
			objectsdb.QueryRow("select Builds from objects WHERE Nme = ?", whoNme).Scan(&bld)
			if bld > 0 {
				objectsdb.Exec("Update objects set Builds = Builds - 1 where Nme = ?", whoNme)
			}
		}
	}
	// // fmt.Println("***points update:hitter:", hitter, hitterSide, "who:", whoNme, whoSide, Objtype, "shared size:",sharedPhasSize, " HIT SIZE:",hitSize)
	doPointsUpdate(hitter, hitterSide, whoNme, whoSide, Objtype, sharedPhasSize, pointsdb, usersdb)

	return
}

//
//  If damage is 300 or more units, the device is inoperative - turn off devices that toggle
//
func dodamagecheck(username string, conn *Conn, objectsdb *sql.DB) {
	var Shld int
	var Cmp int
	var Radio int
	var TractorOn int
	var ShldDam int
	var CmpDam int
	var RadioDam int
	var TractorDam int
	var OutputLen int

	_ = objectsdb.QueryRow("select Shld, ShldDam, Cmp, CmpDam, Radio, RadioDam, TractorOn, TractorDam, OutputLen  from objects WHERE Nme = ?", username).Scan(&Shld, &ShldDam, &Cmp, &CmpDam, &Radio, &RadioDam, &TractorOn, &TractorDam, &OutputLen)

	switch true {
	case ShldDam >= MaxDam && Shld == 1:
		if OutputLen == OutLenLong {
			sendln(conn, "Shields has been knocked offline due to critical damage, Captain!")
		} else {
			if OutputLen == OutLenMed {
				sendln(conn, "Shields offline due to critical damage, Captain!")
			} else {
				sendln(conn, "Shields knocked offline, Captain!")
			}
		}
		objUpdateShldUp(Off, username, objectsdb, conn)

	case CmpDam >= MaxDam && Cmp == 1:
		if OutputLen == OutLenLong {
			sendln(conn, "Computer has been knocked offline due to critical damage, Captain!")
		} else {
			if OutputLen == OutLenMed {
				sendln(conn, "Computer offline due to critical damage, Captain!")
			} else {
				sendln(conn, "Computer knocked offline, Captain!")
			}
		}
		objUpdateCmpUp(0, username, objectsdb)

	case RadioDam >= MaxDam && Radio == 1:
		if OutputLen == OutLenLong {
			sendln(conn, "Radio has been knocked offline due to critical damage, Captain!")
		} else {
			if OutputLen == OutLenMed {
				sendln(conn, "Radio offline due to critical damage, Captain!")
			} else {
				sendln(conn, "Radio knocked offline, Captain!")
			}
		}
		objUpdateRadioUp(0, username, objectsdb)

	case TractorDam >= MaxDam && TractorOn == 1:
		sendln(conn, "Tractor beam has been knocked offline due to critical damage, Captain!")
		objUpdateTractorUp(0, username, objectsdb)
	}
}

//
// Function to update the objects database - prompt damage)
//
func objUpdatePrompt(username string, prompt string, objectsdb *sql.DB) {
	// mu.Lock()
	objectsdb.Exec("UPDATE objects set Prompt = ?  WHERE Nme = ?", prompt, username)
	// fmt.Println("updateprompt to objects:",err, prompt, username)
	// mu.Unlock()
	return
}

//
// Function to update the users db - password)
//
func objUpdatePass(username string, password string, usersdb *sql.DB) {
	usersdb.Exec("UPDATE users set pswd = ? where name = ?", GetMD5Hash(password), username)
	return
}

//
// Function to update the objects database - phaser damage)
//
func objUpdatePhasDam(username string, damage int, objectsdb *sql.DB) {
	var PhasDam int
	_ = objectsdb.QueryRow("select PhasDam from objects WHERE Nme = ?", username).Scan(&PhasDam)
	PhasDam = PhasDam + damage
	objectsdb.Exec("UPDATE objects set PhasDam = ?  WHERE Nme = ?", PhasDam, username)
	return
}

//
// Function to update the objects database - impulse damage)
//
func objUpdateImpDam(username string, damage int, objectsdb *sql.DB) {
	var ImpEngDam int
	_ = objectsdb.QueryRow("select ImpEngDam from objects WHERE Nme = ?", username).Scan(&ImpEngDam)
	ImpEngDam = ImpEngDam + damage
	objectsdb.Exec("UPDATE objects set ImpEngDam = ?  WHERE Nme = ?", ImpEngDam, username)
	return
}

//
// Function to update the objects database - engine damage)
//
func objUpdateEngineDam(username string, damage int, objectsdb *sql.DB) {
	var WarpEngDam int
	_ = objectsdb.QueryRow("select WarpEngDam from objects WHERE Nme = ?", username).Scan(&WarpEngDam)
	WarpEngDam = WarpEngDam + damage
	objectsdb.Exec("UPDATE objects set WarpEngDam = ?  WHERE Nme = ?", WarpEngDam, username)
	return
}

//
// Function to update the objects database - stardate (both game and lifetime)
//
func objUpdateStarDate(username string, objectsdb *sql.DB, usersdb *sql.DB, pointsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects SET StarDate = Stardate + 1  WHERE Nme = ?", username)
	usersdb.Exec("UPDATE users SET NumOfStarDates = NumOfStarDates + 1 WHERE name = ?", username)
	pointsdb.Exec("UPDATE points set NumOfStarDates = NumOfStarDates + 1 WHERE Nme = ?", username)
	return
}

//
// Function to update the objects database - ship status
//
func objUpdateStat(Stat int, username string, objectsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects set Stat = ?  WHERE Nme = ?", Stat, username)
	return
}

//
// Function to update the objects database - active flag
//
func objUpdateActv(Actv int, username string, objectsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects set Actv = ?  WHERE Nme = ?", Actv, username)
	return
}

//
// Function to update the objects database - shield energy
//
func objUpdateShld(Shld int, username string, objectsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects set Shld = ?  WHERE Nme = ?", Shld, username)
	return
}

//
// Function to update the objects database - shield status (up/down)
//
func objUpdateShldUp(ShldUp int, username string, objectsdb *sql.DB, conn *Conn) bool {
	var ShldDam int
	if ShldUp == On {
		objectsdb.QueryRow("select ShldDam  from objects WHERE Nme = ?", username).Scan(&ShldDam)
		if ShldDam >= MaxDam {
			sendln(conn, unableRaiseShld)
			return false
		}
	}
	objectsdb.Exec("UPDATE objects set Shldup = ?  WHERE Nme = ?", ShldUp, username)
	return true
}

//
// Function to update the objects database - computer status (up/down)
//
func objUpdateCmpUp(Cmp int, username string, objectsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects set Cmp = ?  WHERE Nme = ?", Cmp, username)
	return
}

//
// Function to update the objects database - radio status (up/down)
//
func objUpdateRadioUp(Radio int, username string, objectsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects set Radio = ?  WHERE Nme = ?", Radio, username)
	return
}

//
// Function to update the objects database - tractor status (on/off)
//
func objUpdateTractorUp(TractorOn int, username string, objectsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects set TractorOn = ?  WHERE Nme = ?", TractorOn, username)
	return
}

//
// Function to update the objects database - Ship Energy
//
func objUpdateShipEnergy(ShipEnergy int, username string, objectsdb *sql.DB) {
	objectsdb.Exec("UPDATE objects set ShipEnergy = ? WHERE Nme = ?", ShipEnergy, username)
	if ShipEnergy > 1000 {
		objUpdateStat(StatG, username, objectsdb)
	} else {
		objUpdateStat(StatY, username, objectsdb)
	}
	return
}

//
func checkErr(err error) {
	if err != nil {
		if err == io.EOF {
			return
		}
	}
	return
}

//
// Do computed - Return closest object matching obj (abbreviation)
//
func docomputed(username string, obj string, objectsdb *sql.DB) (int, int, bool) {
	var myLocx int
	var myLocy int
	var clLocx int
	var clLocy int
	var mySide int
	var rslt bool

	// Look up location for user (shooter)
	_ = objectsdb.QueryRow("select Locx, Locy, Side from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &mySide)

	switch true {
	case testCommandMatch(NameCoalition, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideCoalition, objectsdb)

	case testCommandMatch(NameEmpire, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideEmpire, objectsdb)

	case testCommandMatch(NameNeutral, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideNeutral, objectsdb)

	case testCommandMatch(NameArcheron, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideArcheron, objectsdb)

	case testCommandMatch(NameFriendlies, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, mySide, objectsdb)

	case (testCommandMatch(NameEnemy, obj, lenOfString1) || testCommandMatch(NameTargets, obj, lenOfString1)):
		if mySide == SideEmpire {
			clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideCoalition, objectsdb)
		} else {
			clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideEmpire, objectsdb)
		}
	case testCommandMatch(NameShips, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeShip, 0, objectsdb)

	case testCommandMatch(NamePlanets, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypePlanet, 0, objectsdb)

	case testCommandMatch(NameStars, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeStar, 0, objectsdb)

	case testCommandMatch(NameBases, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeBase, 0, objectsdb)

	case testCommandMatch(NameBH, obj, lenOfString1):
		clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeBH, 0, objectsdb)
		// if it's a obj, return it's addr!
	default:
		_ = objectsdb.QueryRow("select Locx, Locy from objects WHERE Nme = ?", obj).Scan(&clLocx, &clLocy)
		rslt = true
	}
	return clLocx, clLocy, rslt
}

//
// Function to return string indicator of side
//
func strSide(side int) string {
	if side == SideEmpire {
		return "-"
	}
	if side == SideNeutral {
		return " "
	}
	if side == SideCoalition {
		return "+"
	}
	if side == SideArcheron {
		return " "
	}
	// fmt.Println("got a invalid side:", side)
	return "?"
}

//
// for some commands short out ptu for nmeaddr makes no sense, so short = med
//
func shrtTomed(OutputLen int) int {
	if OutputLen == OutLenSh {
		return OutLenMed
	}
	return OutputLen
}

//
// This function returns a string with the ioFMT format of a user's name/addr and shield strength relative to another location
//  to the relative locations iofmt & OutputLen (whew)
//
func nmeaddr(username string, rellocx int, rellocy int, relioFMT int, relOutputLen int, objectsdb *sql.DB) string {
	var userLocx int
	var userLocy int
	var Shld int
	var ShldUp int
	var Side int
	var Objtype int
	var Builds int
	var msgToReturn string
	var j string
	//
	// Step 1: get Username's address
	//
	objectsdb.QueryRow("select Locx, Locy, Shld, ShldUp, Side, Objtype, Builds  from objects where Nme=?", username).Scan(&userLocx, &userLocy, &Shld, &ShldUp, &Side, &Objtype, &Builds)

	//
	// Step 2: get format to output and do so
	//
	switch relioFMT {
	case IOFmtAbs:
		switch relOutputLen {
		case OutLenSh:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + "-"
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + "+"
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			switch Objtype {
			case TypeShip:
				msgToReturn = msgToReturn + username
				break
			case TypePlanet:
				msgToReturn = msgToReturn + "@" + " (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + "* " + username
				break

			case TypeBase:
				msgToReturn = msgToReturn + "B " + username
				break

			case TypeBH:
				msgToReturn = msgToReturn + "B-H " + username
				break
			}
			msgToReturn = msgToReturn + ", " + j + "%"
			return msgToReturn

		case OutLenMed:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + NameEmpire + " "
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + NameCoalition + " "
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + NameNeutral + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + NameArcheron + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			//			msgToReturn = strSide(Side)
			switch Objtype {
			case TypeShip:
				msgToReturn = msgToReturn + username
				break

			case TypePlanet:
				msgToReturn = msgToReturn + "@ (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + username
				break

			case TypeBase:
				msgToReturn = msgToReturn + username
				break

			case TypeBH:
				msgToReturn = msgToReturn + username
				break
			}
			msgToReturn = msgToReturn + " @" + strconv.Itoa(userLocx) + "-" + strconv.Itoa(userLocy) + ", " + j + "%"
			return msgToReturn

		case OutLenLong:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + NameEmpire + " "
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + NameCoalition + " "
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + NameNeutral + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + NameArcheron + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)

			switch Objtype {
			/*			case TypeShip:
						if Side == SideArcheron {
							msgToReturn = msgToReturn + NameArcheron + " (" + username + ")"
						} else {
							msgToReturn = msgToReturn + username
						}
						break
			*/
			case TypePlanet:
				msgToReturn = msgToReturn + NamePlanetL + " (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + NameStarL + " (" + username + ")"
				break

			case TypeBase:
				msgToReturn = msgToReturn + NameBaseL + " (" + username + ")"
				break

			case TypeBH:
				msgToReturn = msgToReturn + NameBHL + " (" + username + ")"
				break
			}
			msgToReturn = msgToReturn + " @" + strconv.Itoa(userLocx) + "-" + strconv.Itoa(userLocy) + ", " + j + "%"
			return msgToReturn

		}
	case IOFmtRel:
		switch relOutputLen {
		case OutLenSh:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + "-"
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + "+"
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			//			msgToReturn = strSide(Side)
			switch Objtype {
			case TypeShip:
				msgToReturn = msgToReturn + username
				break

			case TypePlanet:
				msgToReturn = msgToReturn + "@ (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + "* " + username
				break

			case TypeBase:
				msgToReturn = msgToReturn + "B " + username
				break

			case TypeBH:
				msgToReturn = msgToReturn + "B-H " + username
				break
			}
			msgToReturn = msgToReturn + ", " + j + "%"
			return msgToReturn

		case OutLenMed:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + NameEmpire + " "
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + NameCoalition + " "
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + NameNeutral + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + NameArcheron + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			//			msgToReturn = strSide(Side)
			switch Objtype {
			case TypeShip:
				msgToReturn = msgToReturn + username
				break

			case TypePlanet:
				msgToReturn = msgToReturn + "@ (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + username
				break

			case TypeBase:
				msgToReturn = msgToReturn + username
				break

			case TypeBH:
				msgToReturn = msgToReturn + username
				break
			}
			msgToReturn = msgToReturn + " " + strconv.Itoa(userLocx-rellocx) + " " + strconv.Itoa(userLocy-rellocy) + ", " + j + "%"
			return msgToReturn

		case OutLenLong:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + NameEmpire + " "
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + NameCoalition + " "
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + NameNeutral + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + NameArcheron + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			//			msgToReturn = strSide(Side)
			switch Objtype {
			case TypeShip:
				/*				if Side == SideArcheron {
									msgToReturn = msgToReturn + NameArcheron + " (" + username + ")"
								} else {
									msgToReturn = msgToReturn + username
								}
								break
				*/
			case TypePlanet:
				msgToReturn = msgToReturn + NamePlanetL + "(" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + NameStarL + " (" + username + ")"
				break

			case TypeBase:
				msgToReturn = msgToReturn + NameBaseL + " (" + username + ")"
				break

			case TypeBH:
				msgToReturn = msgToReturn + NameBHL + " (" + username + ")"
				break
			}
			msgToReturn = msgToReturn + " " + strconv.Itoa(userLocx-rellocx) + " " + strconv.Itoa(userLocy-rellocy) + ", " + j + "%"
			return msgToReturn

		}
	case IOFmtBoth:
		switch relOutputLen {
		case OutLenSh:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + "-"
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + "+"
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			//			msgToReturn = strSide(Side)
			switch Objtype {
			case TypeShip:
				msgToReturn = msgToReturn + username
				break

			case TypePlanet:
				msgToReturn = msgToReturn + "@ (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + "* " + username
				break

			case TypeBase:
				msgToReturn = msgToReturn + "B " + username
				break

			case TypeBH:
				msgToReturn = msgToReturn + "B-H " + username
				break
			}
			msgToReturn = msgToReturn + ", " + j + "%"
			return msgToReturn

		case OutLenMed:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + NameEmpire + " "
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + NameCoalition + " "
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + NameNeutral + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + NameArcheron + " "
			}
			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			//			msgToReturn = strSide(Side)
			switch Objtype {
			case TypeShip:
				msgToReturn = msgToReturn + username
				break

			case TypePlanet:
				msgToReturn = msgToReturn + "@ (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + username
				break

			case TypeBase:
				msgToReturn = msgToReturn + username
				break

			case TypeBH:
				msgToReturn = msgToReturn + username
				break
			}
			msgToReturn = msgToReturn + " @" + strconv.Itoa(userLocx) + "-" + strconv.Itoa(userLocy) + ", " + strconv.Itoa(userLocx-rellocx) + " " + strconv.Itoa(userLocy-rellocy) + ", " + j + "%"
			return msgToReturn

		case OutLenLong:
			//	Show what side it's on (always short format)
			if Side == SideEmpire {
				msgToReturn = msgToReturn + NameEmpire + " "
			}
			if Side == SideCoalition {
				msgToReturn = msgToReturn + NameCoalition + " "
			}
			if Side == SideNeutral {
				msgToReturn = msgToReturn + NameNeutral + " "
			}
			if Side == SideArcheron {
				msgToReturn = msgToReturn + NameArcheron + " "
			}

			pct := calcShields(Shld)
			if ShldUp == On {
				j = "+"
			} else {
				j = "-"
			}

			j = j + fmt.Sprintf("%.2f", pct)
			//			msgToReturn = strSide(Side)
			switch Objtype {
			case TypeShip:
				msgToReturn = msgToReturn + username
				break

			case TypePlanet:
				msgToReturn = msgToReturn + NamePlanetL + " (" + username + ")(" + strconv.Itoa(Builds) + ")"
				break

			case TypeStar:
				msgToReturn = msgToReturn + NameStarL + " (" + username + ")"
				break

			case TypeBase:
				msgToReturn = msgToReturn + NameBaseL + " (" + username + ")"
				break

			case TypeBH:
				msgToReturn = msgToReturn + NameBHL + " (" + username + ")"
				break
			}
			msgToReturn = msgToReturn + " @" + strconv.Itoa(userLocx) + "-" + strconv.Itoa(userLocy) + ", " + strconv.Itoa(userLocx-rellocx) + " " + strconv.Itoa(userLocy-rellocy) + ", " + j + "%"
			return msgToReturn

		}

		msgToReturn = "error on nmeaddr"
	}
	return msgToReturn
}

//
// Send message to multiple folk:
//  Everyone in a location range surrounding an absolute address!
//  Everyone in the game
//  Message type determines this as well as everything else
//
//
//
// Call this with 3 messages: Short, medium, and long formats
// From decwar:
//      This routine processes the information stored in the hit queue,
//      printing out the text produced during battles, primarily.
//      This information is stored in the hit queue using MAKHIT, and
//      is retrieved by OUTHIT using GETHIT.  The 'type' of message is
//      coded into the variable IWHAT:  1=phaser hit, 2=torpedo hit,
//      3=torpedo deflection, 4=torpedo miss, 5=tordedo into black hole,
//      6=star unaffected by torpedo, 7=star goes nova, 8=star damages
//      someone, 9=galaxy-wide base request for assistance, 10=galaxy-
//      wide report of base destroyed, 11=romulan detected message,
//      12=ship-to-ship energy transfer, 13=Tractor beam activated,
//      14=Tractor beam broken, 15=torpedo neutralized.
//		phasHit = size of phaser hit
//
//    Username @22-31 +4,+2 swallowed by black hole          =   locx, locy, message converted to your format for maxrange around locx,locy
//
func notify(username string, conn *Conn, userLocx int, userLocy int, messageType int, impactObjNme string, impactedObjLocx int, impactedObjLocy int, objectsdb *sql.DB, phasHit int, numTorp int) {
	var msgToSend string
	var headerlow int
	var headerhigh int
	var rowlow int
	var rowhigh int
	var Nme string
	var Locx int
	var Locy int
	var IOFmt int
	var OutputLen int
	var Rdio int
	var qry string
	var whatHappened string
	//
	// Get the range to send to
	// 2 types of messages: send to everyone and send to everyone within hearing distance of username
	//
	// Below is send in hearing distance:
	if messageType == msgStar || messageType == msgStarNova || messageType == msgTorpNeu || messageType == msgGulp || messageType == msgGulpTorp || messageType == msgDisp || messageType == msgObjSwallowedBH || messageType == msgObjDied || messageType == msgPhas || messageType == msgTor || messageType == msgTorpMisfire || messageType == msgTorpHit || messageType == msgDestroyed {
		//
		// Compute range in absolute address format
		//
		headerlow = userLocy - MaxScanRng
		headerhigh = userLocy + MaxScanRng
		if headerlow < 0 {
			headerlow = 0
		}
		if headerhigh > Vmax {
			headerhigh = Vmax
		}
		rowlow = userLocx - MaxScanRng
		rowhigh = userLocx + MaxScanRng
		if rowlow < 0 {
			rowlow = 0
		}
		if rowhigh > Hmax {
			rowhigh = Hmax
		}
	} else { // send to everyone
		headerlow = 0
		headerhigh = Vmax
		rowlow = 0
		rowhigh = Hmax
	}
	//
	// now we have determined the range for the message, get the list of ships to send to and send to them
	//
	qry = "select Nme, Locx, Locy, ioFMT, OutputLen, Radio from objects where (Locx between " + strconv.Itoa(rowlow) + " and " + strconv.Itoa(rowhigh) + ") and (Locy between " + strconv.Itoa(headerlow) + " and " + strconv.Itoa(headerhigh) + ") and Objtype = " + strconv.Itoa(TypeShip)
	rows, err := objectsdb.Query(qry)
	if err == nil {
		for rows.Next() {
			rows.Scan(&Nme, &Locx, &Locy, &IOFmt, &OutputLen, &Rdio)
			// Radio must be up!
			if Rdio == On {
				switch OutputLen {
				case OutLenLong:
					switch messageType {
					case msgObjSwallowedBH:
						whatHappened = " swallowed by "
						break

					case msgObjDied:
						whatHappened = " has left the game "
						break

					case msgPhas:
						whatHappened = " makes "
						whatHappened = whatHappened + strconv.Itoa(phasHit)
						whatHappened = whatHappened + " unit phaser hit on "
						break

					case msgStar:
						whatHappened = " makes "
						whatHappened = whatHappened + strconv.Itoa(phasHit)
						whatHappened = whatHappened + " unit hit on "
						break

					case msgTor:
						whatHappened = " shoots torpedo at "
						break

					case msgTorpMisfire:
						whatHappened = " Torpedo MISFIRES!!"
						break

					case msgTorpHit:
						whatHappened = " makes "
						whatHappened = whatHappened + strconv.Itoa(phasHit)
						whatHappened = whatHappened + " unit torpedo #"
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " hit on "
						break

					case msgTorpNeu:
						whatHappened = " torpedo "
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " neutralized by friendly object "
						break

					case msgTorpMiss:
						whatHappened = " torpedo "
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " MISS!!!! "
						break

					case msgStarNova:
						whatHappened = " destroyed by nova"
						break

					case msgGulp:
						whatHappened = bhGulp
						break

					case msgGulpTorp:
						whatHappened = gulpTorp
						break

					case msgDisp:
						whatHappened = " displaced to "
						break

					case msgEndGame:
						whatHappened = endGame
						break

					case msgEndCoalition:
						whatHappened = endGameWhoC
						break

					case msgEndEmpire:
						whatHappened = endGameWhoE
						break
					}

				case OutLenMed:
					switch messageType {
					case msgObjSwallowedBH:
						whatHappened = " quaffed by "
						break

					case msgObjDied:
						whatHappened = " has died "
						break

					case msgPhas:
						whatHappened = " " + strconv.Itoa(phasHit)
						whatHappened = whatHappened + " unit phaser on "
						break

					case msgStar:
						whatHappened = " " + strconv.Itoa(phasHit)
						whatHappened = whatHappened + " unit blast on "
						break

					case msgTor:
						whatHappened = " torpedos at "
						break

					case msgTorpMisfire:
						whatHappened = " Torpedo MISFIRES!!"
						break

					case msgTorpHit:
						whatHappened = " torpedo #"
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " hit on "
						break

					case msgTorpNeu:
						whatHappened = " torpedo "
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " neutralized "
						break

					case msgTorpMiss:
						whatHappened = " T"
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " MISS!!!! "
						break

					case msgStarNova:
						whatHappened = " destroyed by nova"
						break

					case msgGulp:
						whatHappened = bhGulp
						break

					case msgGulpTorp:
						whatHappened = gulpTorp
						break

					case msgDisp:
						whatHappened = " displaced to "
						break

					case msgEndGame:
						whatHappened = endGame
						break

					case msgEndCoalition:
						whatHappened = endGameWhoC
						break

					case msgEndEmpire:
						whatHappened = endGameWhoE
						break

					}

				case OutLenSh:
					switch messageType {
					case msgObjSwallowedBH:
						whatHappened = " gulped by "
						break

					case msgObjDied:
						whatHappened = " died "
						break

					case msgPhas:
						whatHappened = " " + strconv.Itoa(phasHit)
						whatHappened = whatHappened + " P "
						break

					case msgStar:
						whatHappened = " " + strconv.Itoa(phasHit)
						whatHappened = whatHappened + " N "
						break

					case msgTor:
						whatHappened = " T "
						break

					case msgTorpMisfire:
						whatHappened = " Torpedo MISFIRES!!"
						break

					case msgTorpHit:
						whatHappened = " T #"
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " hit on "
						break

					case msgTorpNeu:
						whatHappened = " T "
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " neutralized "
						break

					case msgTorpMiss:
						whatHappened = " T"
						whatHappened = whatHappened + strconv.Itoa(numTorp)
						whatHappened = whatHappened + " MISS!!!! "
						break

					case msgStarNova:
						whatHappened = " destroyed by nova"
						break

					case msgGulp:
						whatHappened = bhGulp
						break

					case msgGulpTorp:
						whatHappened = gulpTorp
						break

					case msgDisp:
						whatHappened = " displaced to "
						break

					case msgEndGame:
						whatHappened = endGame
						break

					case msgEndCoalition:
						whatHappened = endGameWhoC
						break

					case msgEndEmpire:
						whatHappened = endGameWhoE
						break
					}
				}

				// build message that is unique to that receiver!
				switch messageType {
				case msgObjSwallowedBH:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgObjDied:
					msgToSend = username + whatHappened
					break

				case msgPhas:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgStar:
					msgToSend = username + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgTor:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgTorpMisfire:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened
					break

				case msgBuildBase:
					msgToSend = username + BPB

				case msgDestroyed:
					msgToSend = username + Destroyed

				case msgTorpHit:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgTorpNeu:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgTorpMiss:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgStarNova:
					msgToSend = nmeaddr(username, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened
					//					fmt.Println("**in shout notification for star nova: user:", username, " lockx:", Locx, " Locy:", Locy, " whathappened:", whatHappened, " nme:", Nme)
					break

				case msgGulp:
					msgToSend = nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened + username
					break

				case msgGulpTorp:
					msgToSend = nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb) + whatHappened
					break

				case msgDisp:
					msgToSend = username + whatHappened + nmeaddr(impactObjNme, Locx, Locy, IOFmt, OutputLen, objectsdb)
					break

				case msgEndGame:
					msgToSend = whatHappened
					break

				case msgEndCoalition:
					msgToSend = whatHappened
					break

				case msgEndEmpire:
					msgToSend = whatHappened
					break
				}
				//				If the connection is dead, don't send a message!
				if Conmap[Nme].Connection != nil {
					sendln(Conmap[Nme].Connection.(*Conn), msgToSend)
				}
			}
		}
	}
	rows.Close()
}

//
//  See if the command/parameter input matches the command in question. Returns True/False
//
func testCommandMatch(cmd1 string, inputstring1 string, lenOfString int) bool {
	// Pull out any special characters - everything should be lower case
	cmd := strings.ToLower(strings.Trim(cmd1, ctlsp))
	inputstring := strings.ToLower(strings.Trim(inputstring1, ctlsp))
	// All commands & parms are 2 char or more ... period
	if (len(inputstring) > len(cmd)) || (len(inputstring) < lenOfString) {
		return false
	}
	if cmd[0:len(inputstring)] == inputstring {
		return true
	} else {
		return false
	}
}

//
// Generate random password
//
func doRandomPassword() string {
	maxPasswordLength := 8
	password := make([]byte, maxPasswordLength)
	alphabet := "abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	for i := 0; i < len(password); i++ {
		password[i] = alphabet[rand.Int()%len(alphabet)]
	}
	return string(password)
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//
// Calc incremental moves
// Used for moves and torps, calculate the # moves needed, along with the incremental x, y for each move
//
func doCalcIncMove(relx int, rely int) (nummoves int, incx float32, incy float32) {
	// Which move is bigger, x or y?
	if Abs(relx) > Abs(rely) {
		nummoves = Abs(relx)
	} else {
		nummoves = Abs(rely)
	}
	// Calc incremental move for each
	incx = float32(relx) / float32(nummoves)
	incy = float32(rely) / float32(nummoves)
	return nummoves, incx, incy
}

//
// Do the help command
//
func doHelp(invalidcommand bool, conn *Conn, nme string, syntx string, fctn string, optns string) bool {
	send(conn, "NAME:")
	send(conn, "\t")
	sendln(conn, nme)
	send(conn, "SYNTAX:")
	send(conn, "\t")
	sendln(conn, syntx)
	send(conn, "FUNCTION:")
	send(conn, "\t")
	sendln(conn, fctn)
	send(conn, "OPTIONS:")
	send(conn, "\t")
	sendln(conn, optns)
	invalidcommand = false
	return invalidcommand
}

//
// Do the help examples (multiple lines sent 1 at a time)
//
func doExample(invalidcommand bool, conn *Conn, example string, syntx string, fctn string, optns string) bool {
	send(conn, "Example:")
	send(conn, "\t")
	sendln(conn, example)
	invalidcommand = false
	return invalidcommand
}

//
// Do the default pregame help command
//
func dodefPregameHelp(invalidcommand bool, conn *Conn) bool {
	sendln(conn, "Help command:")
	sendln(conn, "Syntax: help [command]")
	sendln(conn, "Lists or describe the legal commands.")
	sendln(conn, "Commands available:")
	sendln(conn, "\tActivate \tGripe \t\tHelp \t\tLOGIn")
	sendln(conn, "\tLOGOff \t\tNews \t\tPoints \t\tREGister")
	sendln(conn, "\tSUmmary \tTime	\tUsers \t\tQuit")
	sendln(conn, "\tREcover\n")
	sendln(conn, "Command description syntax  (DO NOT INCLUDE [,], |, <, >, { OR } WHEN TYPING COMMANDS!):")
	sendln(conn, "\t*   Square brackets for optional parameters: help [list]")
	sendln(conn, "\t*   Angle brackets for required parameters: register <username> <email addr>")
	sendln(conn, "\t*   Ellipses for repeated items: send <user> {word ... word}")
	sendln(conn, "\t*   Vertical bars for choice of items: list {enemy | friendly}\n")
	sendln(conn, "Commands can be stacked on one line by seperating with a /")
	sendln(conn, "\t    Example:   help/news/time")
	sendln(conn, "You can repeat the previous command line by pressing the enter key on a blank line.\n")
	sendln(conn, "USE THE ACTIVATE COMMAND TO BE ASSIGNED A USERNAME AND ENTER THE GAME AS A GUEST!\n")
	sendln(conn, "Use the Register to register a new username.\n ")
	sendln(conn, "Use Login to log into the server.")
	return false
}

//
//
//
// Do the zero script
//
func doZero(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Zero string
	//
	// Run the zero command!!!
	//
	usersdb.QueryRow("Select Zero from users where name = ?", username).Scan(&Zero)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Zero, ctlonly), "/")
	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Zero, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the One script
//
func doOne(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var One string
	//
	// Run the One command!!!
	//
	usersdb.QueryRow("Select One from users where name = ?", username).Scan(&One)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(One, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], One, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Two script
//
func doTwo(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Two string
	//
	// Run the One command!!!
	//
	usersdb.QueryRow("Select Two from users where name = ?", username).Scan(&Two)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Two, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Two, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Three script
//
func doThree(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Three string
	//
	// Run the Three command!!!
	//
	usersdb.QueryRow("Select Three from users where name = ?", username).Scan(&Three)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Three, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Three, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Four script
//
func doFour(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Four string
	//
	// Run the Four command!!!
	//
	usersdb.QueryRow("Select Four from users where name = ?", username).Scan(&Four)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Four, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Four, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Five script
//
func doFive(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Five string
	//
	// Run the Five command!!!
	//
	usersdb.QueryRow("Select Five from users where name = ?", username).Scan(&Five)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Five, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Five, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Six script
//
func doSix(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Six string
	//
	// Run the Six command!!!
	//
	usersdb.QueryRow("Select Six from users where name = ?", username).Scan(&Six)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Six, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Six, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Seven script
//
func doSeven(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Seven string
	//
	// Run the Seven command!!!
	//
	usersdb.QueryRow("Select Seven from users where name = ?", username).Scan(&Seven)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Seven, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Seven, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Eight script
//
func doEight(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Eight string
	//
	// Run the Seven command!!!
	//
	usersdb.QueryRow("Select Eight from users where name = ?", username).Scan(&Eight)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Eight, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Eight, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the Nine script
//
func doNine(username string, invalidcommand bool, conn *Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var Nine string
	//
	// Run the Seven command!!!
	//
	usersdb.QueryRow("Select Nine from users where name = ?", username).Scan(&Nine)

	//
	//Parse commands - seperated by /
	//
	comnd1 := strings.Split(strings.Trim(Nine, ctlonly), "/")

	// Don't allow script looping
	for i := range comnd1 {
		//
		// not allowed to stack activate command
		//
		_, err := strconv.Atoi(comnd1[i])
		if err != nil {
			username = processCommand(comnd1[i], Nine, conn, username, err, false, usersdb, objectsdb, pointsdb)
		} else {
			sendln(conn, "You cannot call scripts from scripts, captian!")
		}
	}
	return false
}

//
// Do the default help command
//
func dodefHelp(invalidcommand bool, conn *Conn) bool {
	sendln(conn, "Help command:")
	sendln(conn, "Syntax: help [command]")
	sendln(conn, "Lists or describe the legal commands.")
	sendln(conn, "Commands available:")
	send(conn, "\t\tADministrator")
	send(conn, "\tBAses ")
	send(conn, "\t\tBUild")
	sendln(conn, "\t\tCApture")
	send(conn, "\t\tDAmages")
	send(conn, "\t\tDOck")
	send(conn, "\t\tENergy")
	sendln(conn, "\t\tGRipe")
	send(conn, "\t\tHElp")
	send(conn, "\t\tIMpulse")
	send(conn, "\t\tLIst")
	sendln(conn, "\t\tMOve")
	send(conn, "\t\tNEws")
	send(conn, "\t\tPHasers")
	send(conn, "\t\tPLanets")
	sendln(conn, "\t\tPOints")
	send(conn, "\t\tQUit")
	send(conn, "\t\tRAdio")
	send(conn, "\t\tREpair")
	sendln(conn, "\t\tSCan")
	send(conn, "\t\tSET")
	send(conn, "\t\tSHields")
	send(conn, "\t\tSRscan")
	sendln(conn, "\t\tSTatus")
	send(conn, "\t\tSUmmary")
	send(conn, "\t\tTArgets")
	send(conn, "\t\tTEll")
	sendln(conn, "\t\tTIme")
	send(conn, "\t\tTOrpedoes")
	send(conn, "\tTRactor")
	send(conn, "\t\tTYpe")
	sendln(conn, "\t\tUSers")
	sendln(conn, "\t\t0-9")
	sendln(conn, "Command description syntax  (DO NOT INCLUDE [, |, < OR { WHEN TYPING COMMANDS!):")
	sendln(conn, "\t*   Square brackets for optional parameters: help [list]")
	sendln(conn, "\t*   Angle brackets for required parameters: register <username> <email addr>")
	sendln(conn, "\t*   Ellipses for repeated items: TEll <user> {word ... word}")
	sendln(conn, "\t*   Vertical bars for choice of items: list {enemy | friendly}\n")
	sendln(conn, "Commands can be stacked on one line by seperating with a /")
	sendln(conn, "\t    Example:   help/news/time")
	sendln(conn, "You can repeat the previous command line by pressing the enter key on a blank line.\n")
	sendln(conn, "You can create an intial script set of commands using the set initial command.\n")
	sendln(conn, "You can create scriptable commands using the set # command.\n")
	return false
}

//
// Activate command - add a user to the object database and set flag to get out of pregame
//   INSERT INTO "Objects" (nme, Conndatetime, Locx, Locy) VALUES('test',0,31,1);
func processActivate(comnd string, invalidcommand bool, conn *Conn, username string, netaddr string, objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB) (bool, string) {
	var sde int
	var mySide int
	var myActv int
	var myPrompt string
	var initcmd string
	var err error
	var comnd2 string

// handle guests ie: no login
	if username == "" {
		username = strings.Join(strings.Split(netaddr, ":")[1:], ":")

			// Ok to log them in
		
	Conmap[username] = Constr{username, conn, netaddr}
	}
	
	if username != "" {
		// add to the objects db  **** random setup side if not specified needs to be fixed ******
		for {
			if rand.Intn(2) == 1 {
				sde = SideCoalition
			} else {
				sde = SideEmpire
			}
			// returning customer?
			err := objectsdb.QueryRow("select Actv, Side, Prompt  from objects where Nme=?", username).Scan(&myActv, &mySide, &myPrompt)
			if err == nil {
				objUpdateActv(On, username, objectsdb)
				break
			} else {
				for {
					success := dbAddobjects(usersdb, pointsdb, objectsdb, username, StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), sde, TypeShip, sde, 0, 0, 0, InitPhoTor, 0, 0,
						InitShield, On, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, IOFmtBoth, OutLenLong, 0, "Decwars", "", "")
					// fmt.Println("success:", success)
					if success == true {
						break
					}
				}
			}
		}

		sendln(conn, "Entering decwars...")
		//
		objectsdb.QueryRow("select Actv from objects where Nme=?", username).Scan(&myActv)
		//
		// Run the initial command!!!
		//
		err = usersdb.QueryRow("Select InitialCommand from users where name = ?", username).Scan(&initcmd)

if err == nil {
		//
		//Parse commands - seperated by /
		//
		comnd1 := strings.Split(strings.Trim(initcmd, ctlonly), "/")
		comnd2 = initcmd

		for i := range comnd1 {
			//
			// not allowed to stack activate command - show your are in initial command!
			//
			username = processCommand(comnd1[i], initcmd, conn, username, err, true, usersdb, objectsdb, pointsdb) //
			//
		}
		comnd2 = strings.Join(comnd1[1:], "/")
		comnd1 = strings.Split(strings.Trim(comnd2, ctlonly), "/")
		return false, username
} else {
	return true, username
}
	} else {  //user is not logged in, create temp username
		sendln(conn, "You can't activate until logged in")
		return false, username
	}
}

//
// Do the default bases command - no parms
//
func dodefBases(invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	var mySide int
	var myLocx int
	var myLocy int
	var myIOFmt int
	var myOutputLen int
	// get my info
	objectsdb.QueryRow("select Side, Locx, Locy, IOFmt, OutputLen from objects where Nme=?", username).Scan(&mySide, &myLocx, &myLocy, &myIOFmt, &myOutputLen)

	qry := "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeBase)
	doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
	doSummaryBases(conn, objectsdb)
	sendln(conn, " ")
	return false
}

//
// Account bounce - kill someone
//  Conmap[whoNme].Connection.(*Conn)
func accountBounce(bounceuser string, comnd string, invalidcommand bool, conn *Conn, netaddr string, usersdb *sql.DB, objectsdb *sql.DB) bool {
	var myActv int
	if bounceuser == "" {
		sendln(conn, "Must specify who to kill")
		return false
	} else {
		//
		// User must be logged on
		//
		objectsdb.QueryRow("select Actv from objects where Nme=?", bounceuser).Scan(&myActv)
		if myActv == On {
			log.Print("Decwars bouncing:", bounceuser, Conmap[bounceuser], Conmap[bounceuser].Connection.(*Conn))
			// Tell them why
			sendln(Conmap[bounceuser].Connection.(*Conn), "Admin has bounced you!")
			// and close the connection
			checkErr(Conmap[bounceuser].Connection.(*Conn).Close())
			return false
		} else {
			sendln(conn, "User is not logged on!")
		}
	}
	return false
}

//
// Account delete - kill an object by name
//  Conmap[whoNme].Connection.(*Conn)
//
func accountDelete(deleteuser string, comnd string, invalidcommand bool, conn *Conn, netaddr string, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	//	var objSide int
	var objType int
	//	var objNme string
	//	var objLocx int
	//	var objLocy int
	//	var objShld int
	//	var objShldUp int
	//	var objWarpEngDam int
	//	var objImpEngDam int
	//	var objPhoTorDam int
	//	var objPhasDam int
	//	var objShldDam int
	//	var objCmpDam int
	//	var objLifeSupDam int
	//	var objRadioDam int
	//	var objTractorDam int
	//	var objShipDam int
	//	var hitterNme string
	//	var hitterLocx int
	//	var hitterLocy int

	//	var hitterSide int
	//  var err error
	if deleteuser == "" {
		sendln(conn, "Must specify object to delete")
		return false
	} else {
		//
		// Only delete planets, stars, bases, archerons, black holes
		//
		objectsdb.QueryRow("select Objtype from objects where Nme=?", deleteuser).Scan(&objType)

		//known bad err = objectsdb.QueryRow("select Side, Objtype, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam from objects where Nme=?", deleteuser).Scan(&objSide, &objType, &objLocx, &objLocy, &objShld, &objShldUp, &objWarpEngDam, &objImpEngDam, &objPhoTorDam, &objPhasDam, &objShldDam, &objCmpDam, &objLifeSupDam, &objRadioDam, &objTractorDam, &objShipDam)

		//// fmt.Println("err=",err,objNme, objType, objSide, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objType, 10000, objectsdb, nil, torpHit, hitterNme, hitterLocx, hitterLocy, hitterSide)

		if objType == TypePlanet || objType == TypeStar || objType == TypeBase || objType == TypeBH || objType == TypeShip {
			log.Print("Decwars deleting:", deleteuser)
			// Perform the deletion
			dbDelobjects(nil, objectsdb, deleteuser, pointsdb, usersdb)

			//known  doHit(objNme, objSide, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objType, 10000, objectsdb, nil, torpHit, hitterNme, hitterLocx, hitterLocy, hitterSide, pointsdb, usersdb)
			return false
		} else {
			sendln(conn, "Invalid object type!")
		}
	}
	return false
}

//
// Account Report - run with "" means for everyone
//
func accountReport(rptuser string, comnd string, invalidcommand bool, conn *Conn, netaddr string, usersdb *sql.DB) bool {
	var qry string
	if rptuser == "" {
		qry = "SELECT name, disabled, RecoveryDateSent, mailAddr, SuperUser FROM users order by name"
	} else {
		qry = "SELECT name, disabled, RecoveryDateSent, mailAddr, SuperUser FROM users where name = " + `"` + rptuser + `";`
	}

	rows, err := usersdb.Query(qry)
	if err == nil {
		for rows.Next() {
			var nme string
			var disbled int
			var rds string
			var mail string
			var su int

			rows.Scan(&nme, &disbled, &rds, &mail, &su)
			send(conn, nme)
			send(conn, "\t")
			if disbled == 1 {
				send(conn, "disabled")
			} else {
				send(conn, "enabled ")
			}
			send(conn, "\t")
			send(conn, rds)
			send(conn, "\t")

			if su == 1 {
				send(conn, "admin\t")
			} else {
				send(conn, "non-admin")
			}
			send(conn, "\t")
			sendln(conn, mail)
		}
	}
	return false
}

//
// Account command
// ADministrator report ALL
//
func dodefAccount(username string, comnd string, invalidcommand bool, conn *Conn, netaddr string, usersdb *sql.DB) bool {
	invalidcommand = processHelp("help admin", invalidcommand, conn)
	return false
}

//
// ADministrator command
// ADministrator <bounce | disable | enable | report | remove | su-on | su-off | bh-off | bh-on | ar-on |ar-off | delete | move | reset | wed | ied | ptd | pbd | sd | cd| lsd | rd | tbd | tsd | dohit | start> <username/object> <Locy Locx> <Amount>
//
func processAdmin(username string, comnd string, invalidcommand bool, conn *Conn, netaddr string, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) bool {
	var su int
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var parameter4 string
	var parameter5 string
	var parameter6 string
	var parameter7 string
	var Locx int
	var Locy int
	var Nme string
	var SeenbyEnemy int
	var MySide int

	var mySide int
	var myLocx int
	var myLocy int
	var myShld int
	var myShldUp int
	var myWarpEngDam int
	var myPhoTorDam int
	var myPhasDam int
	var myShldDam int
	var myCmpDam int
	var myLifeSupDam int
	var myRadioDam int
	var myTractorDam int
	var myShipDam int
	var myObjtype int
	var myImpEngDam int
	var hitamt int

	var hisLocx int
	var hisLocy int
	var hisSide int
	var err error
	var portnum int
	//
//		fmt.Println("starting admin command")

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")

	parameter1 = strings.Trim(prm1[1], ctlsp)

	// must be a super user to use this command
	usersdb.QueryRow("select SuperUser from users where name=?", username).Scan(&su)
//	fmt.Println("got err:", err33, " su:",su)
	if su == 0 {
		sendln(conn, "Not a superuser")
		return false
	} 
//fmt.Println("parameter1=", parameter1, " prm1=", prm1, " len(prm1)", len(prm1))
	// bh-off/bh-on commands
	//fmt.Println("len prm1[1]:",len(prm1))
	// 1 parm provided for report
	if len(prm1) == 2 { //1 parms - must be the name to report or the word report which means do all
		switch true {
	
	case testCommandMatch(NameEndG, parameter1, len(parameter1)):
//fmt.Println("Port:",Portnum," got into endgame")
//		sendln(conn, strconv.Itoa(Portnum))
//		sendln(conn, "got into endgame")
		if Portnum != 1701 {
			sendln(conn, "Ending the game as requested!")
			notify("", nil, 0, 0, msgEndGame, "", 0, 0, objectsdb, 0, 0)
			os.Exit(3)
		} else {
			sendln(conn, "Ending the port 1701 game is disallowed!")
		}
		return false

	case testCommandMatch(NameBHoff, parameter1, len(parameter1)):
		var count int
		objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Objtype = ?", TypeBH).Scan(&count)
		if count != 0 {
			objectsdb.Exec("Delete from objects where Objtype=?", TypeBH)
			sendln(conn, "Black holes have been turned off")
		} else {
			sendln(conn, "Black holes are already off!")
		}
		return false

	case testCommandMatch(NameBHon, parameter1, len(parameter1)):
		var count int
		objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Objtype = ?", TypeBH).Scan(&count)
		if count == 0 {
			// Build black holes
			//		BHsmax := Abs(int(41 * ((rand.NormFloat64() * .025) + 10)))
			for i := 0; i < int(BHsmax); i++ {
				//	for i := 0; i < int(3); i++ { for testing without bh
				success := false
				for {
					success = dbAddobjects(usersdb, pointsdb, objectsdb, "bh"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideNeutral, TypeBH, SideNeutral, 0, 0, 0, InitPhoTor, 0, 0,
						InitShield, Off, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
					if success == true {
						break
					}
				}
			}
			sendln(conn, "Black holes have been turned on")
		} else {
			sendln(conn, "Black holes are already on!")
		}
		return false

		//
		// ar-on and ar-off commands
		//
	case testCommandMatch(NameARon, parameter1, len(parameter1)):
		//
		// Wake the Archeron up
		// Build Archeron ship(s)
		//
		var count int
		objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ?", SideArcheron).Scan(&count)
		if count == 0 {
			for i := 0; i < int(Archeronmax); i++ {
				success := false
				for {
					success = dbAddobjects(usersdb, pointsdb, objectsdb, "a"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideArcheron, TypeShip, SideArcheron, 0, 0, 0, InitPhoTor, 0, 0,
						InitShield, On, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
					if success == true {
						break
					}
				}
			}
			sendln(conn, "Archerons have been turned on")
		} else {
			sendln(conn, "Archerons already are on!")
		}
		return false

	case testCommandMatch(NameARoff, parameter1, len(parameter1)):
		var count int
		objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ?", SideArcheron).Scan(&count)
		if count != 0 {
			objectsdb.Exec("Delete from objects where Side=?", SideArcheron)
			sendln(conn, "Archerons have been turned off")
		} else {
			sendln(conn, "Archerons are already off!")
		}
		return false


	//	fmt.Println("length of prm1:", len(prm1), " and prm1:", prm1)


		case testCommandMatch(NameReport, parameter1, len(parameter1)):
			accountReport("", "", invalidcommand, conn, netaddr, usersdb)
			return false
		}

		accountReport(parameter1, "", invalidcommand, conn, netaddr, usersdb)
		return false
	}

	// 2 parms for bounce delete enable and disable su-on su-off
	if len(prm1) == 3 { //2 parms - must be command and name
		parameter2 = strings.Trim(prm1[2], ctlsp) // must be a username
		// fmt.Println("In 3")
		switch true {
		case testCommandMatch(NameEnable, parameter1, len(parameter1)):
			usersdb.Exec("UPDATE users set disabled = 0 where name=?", parameter2)
			return false

		case testCommandMatch(NameDisable, parameter1, len(parameter1)):
			if testCommandMatch("admin", parameter2, 5) || testCommandMatch("your password here", parameter2, 5) {
				sendln(conn, "Cannot change admin")
				return false
			}

			usersdb.Exec("UPDATE users set disabled = 1 where name=?", parameter2)
			return false

			//
			// Remove the user provided from the db - sqlite won't return an error if non-existant name given! Dont delete hsn or admin
			//
		case testCommandMatch(NameRemove, parameter1, len(parameter1)):
			if testCommandMatch("admin", parameter2, 5) || testCommandMatch("your password here", parameter2, 5) {
				sendln(conn, "Cannot change admin")
				return false
			}
			usersdb.Exec("Delete from users where name=?", parameter2)
			sendln(conn, "Removed user")
			return false

		// Bounce (kill) someone
		case testCommandMatch(NameReset, parameter1, len(parameter1)):
			var email string
			//	reset someone's pw
			parameter2 = strings.Trim(prm1[2], ctlsp)
			usersdb.QueryRow("select mailAddr from users where name=?", parameter2).Scan(&email)
			commnd := "recover " + email

			processRecover(commnd, invalidcommand, conn, conn.Conn.RemoteAddr().String(), usersdb)
			return false

		// Bounce (kill) someone
		case testCommandMatch(NameBounce, parameter1, len(parameter1)):
			//	kill someone
			parameter2 = strings.Trim(prm1[2], ctlsp)
			accountBounce(parameter2, "", invalidcommand, conn, netaddr, usersdb, objectsdb)
			return false

		// report provided with name
		case testCommandMatch(NameReport, parameter1, len(parameter1)):
			parameter2 = strings.Trim(prm1[2], ctlsp)
			accountReport(parameter2, "", invalidcommand, conn, netaddr, usersdb)
			return false

		// Delete command - kill an object
		case testCommandMatch(NameDelete, parameter1, len(parameter1)):
			parameter2 = strings.Trim(prm1[2], ctlsp)
			accountDelete(parameter2, "", invalidcommand, conn, netaddr, usersdb, objectsdb, pointsdb)
			return false

		// turn a player into admin account
		case testCommandMatch(NameSuon, parameter1, len(parameter1)):

			// mu.Lock()
			usersdb.Exec("UPDATE users set SuperUser = 1 where name=?", parameter2)
			// mu.Unlock()
			return false

			//turn off admin account for player
		case testCommandMatch(NameSuoff, parameter1, len(parameter1)):

			if testCommandMatch("admin", parameter2, len(parameter2)) || testCommandMatch("your password here", parameter2, 5) {
				sendln(conn, "Cannot change admin")
				return false
			}
			if testCommandMatch(username, parameter2, len(parameter2)) {
				sendln(conn, "Cannot turn off administrator privs to your account")
				return false
			}

			// mu.Lock()
			usersdb.Exec("UPDATE users set SuperUser = 0 where name=?", parameter2)
			// mu.Unlock()
			return false
		}
	}

	if len(prm1) == 4 { //3 parms - must be command, name, and #
		parameter2 = strings.Trim(prm1[2], ctlsp) // must be a username
		parameter3 = strings.Trim(prm1[3], ctlsp) // must be a #
		//			fmt.Println("In 4 parameter 2:", parameter2, " parameter3:", parameter3)
		switch true {

		//
		// Change warp engine damage amount
		//
		case testCommandMatch(NameWED, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set WarpEngDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change impulse engine damage amount
			//
		case testCommandMatch(NameIED, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set ImpEngDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change photon torpedoes damage amount
			//
		case testCommandMatch(NamePTD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set PhoTorDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change phaser bank damage amount
			//
		case testCommandMatch(NamePBD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set PhasDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change Shield damage amount
			//
		case testCommandMatch(NameSD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set ShldDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change computer damage amount
			//
		case testCommandMatch(NameCD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set CMPDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change life support damage amount
			//
		case testCommandMatch(NameLSD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set LifeSupDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change radio damage amount
			//
		case testCommandMatch(NameRD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set RadioDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change tractor beam damage amount
			//
		case testCommandMatch(NameTBD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set TractorDam = ?  WHERE Nme = ?", amt, parameter2)
			return false

			//
			// Change total ship damage amount
			//
		case testCommandMatch(NameTSD, parameter1, len(parameter1)):
			amt, err := strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return true
			}

			objectsdb.Exec("UPDATE objects set ShipDam = ?  WHERE Nme = ?", amt, parameter2)
			return false
		}
	}

	//
	// Administrative move command - move an object anywhere in the universe
	//

	//fmt.Println("in admin move len(prm1):", len(prm1))

	if len(prm1) == 5 { //3 parms - must be command, name and locy, locx
		switch true {
		case testCommandMatch("move", parameter1, len(parameter2)):
			parameter2 = strings.Trim(prm1[2], ctlsp) // must be a username
			parameter3 = strings.Trim(prm1[3], ctlsp) // must be a locy
			parameter4 = strings.Trim(prm1[4], ctlsp) // must be a locy
			//		fmt.Println("in 5 move got 2:", parameter2, " 3:", parameter3, " 4:", parameter4, "parm1:", parameter1)
			switch true {
			case testCommandMatch("move", parameter1, len(parameter1)):
				pmy, err1 := strconv.Atoi(parameter4)
				pmx, err2 := strconv.Atoi(parameter3)
				if err1 != nil || err2 != nil {
					sendln(conn, invparm)
					return false
				}
				if pmx <= Hmax && pmx >= 0 && pmy <= Vmax && pmy >= 0 {
					// something in the way?

					err := objectsdb.QueryRow("select Nme, Locx, Locy, Side from objects WHERE Locy= ? and Locx = ?", pmy, pmx).Scan(&Nme, &Locy, &Locx, &MySide)
					// fmt.Println(" got move err:",err, "Locy:", Locy, " Locx:",Locx)
					if err != nil {
						SeenbyEnemy = MySide
						objectsdb.Exec("UPDATE objects set Locx = ?, Locy = ?, SeenbyEnemy = ? WHERE Nme = ?", pmx, pmy, SeenbyEnemy, parameter2)
						sendln(conn, "Move completed")
						return false
					} else {
						sendln(conn, "Move blocked by existing object")
						return false
					}
				} else {
					sendln(conn, invparm)
					return false
				}
			}
		}
		return false
	}

	//
	// Admin dohit: adm dohit hsn wolf phaser 500   (dohit from to type amount)
	//
	if len(prm1) == 6 { //3 parms - must be command, name and locy, locx
		parameter2 = strings.Trim(prm1[2], ctlsp) // must be a username
		parameter3 = strings.Trim(prm1[3], ctlsp) // must be a username
		parameter4 = strings.Trim(prm1[4], ctlsp) // must be a type
		parameter5 = strings.Trim(prm1[5], ctlsp) // must be a amount
		switch true {
		case testCommandMatch(NamedoHit, parameter1, len(parameter1)):

			fmt.Println(", ", parameter6)

			//
			// get the info needed for the hit
			//

			hitamt, err = strconv.Atoi(parameter5)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			objectsdb.QueryRow("select Side, Locx, Locy, Shld, ShldUp, WarpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, ImpEngDam from objects where Nme=?", parameter3).Scan(&mySide, &myLocx, &myLocy, &myShld, &myShldUp, &myWarpEngDam, &myPhoTorDam, &myPhasDam, &myShldDam, &myCmpDam, &myLifeSupDam, &myRadioDam, &myTractorDam, &myShipDam, &myObjtype, &myImpEngDam)
			objectsdb.QueryRow("select Side, Locx, Locy from objects where Nme=?", parameter2).Scan(&hisSide, &hisLocx, &hisLocy)

			switch true {
			case testCommandMatch(NameTorp, parameter4, len(parameter4)):
				doHit(parameter3, mySide, myLocx, myLocy, myShld, myShldUp, myWarpEngDam, myImpEngDam, myPhoTorDam, myPhasDam, myShldDam, myCmpDam, myLifeSupDam, myRadioDam, myTractorDam, myShipDam, myObjtype, hitamt, objectsdb, conn, torpHit, parameter2, myLocx, myLocy, mySide, pointsdb, usersdb)
				return false

			case testCommandMatch(NamePhas, parameter4, len(parameter4)):
				doHit(parameter3, mySide, myLocx, myLocy, myShld, myShldUp, myWarpEngDam, myImpEngDam, myPhoTorDam, myPhasDam, myShldDam, myCmpDam, myLifeSupDam, myRadioDam, myTractorDam, myShipDam, myObjtype, hitamt, objectsdb, conn, phasHit, parameter2, myLocx, myLocy, mySide, pointsdb, usersdb)
				return false

			}
			return false
		}
		return false
	}

	//
	// Start command
	// Format ex: ad start Vmax Hmax Starsmax BHsmax Archeronmax Planetsmax Port#
	//
	//the below length check has some type of bug!
//fmt.Println("pre: got into start command! len(prm1):", len(prm1))
	if len(prm1) == 9 { // this changes if you add more parms
		if testCommandMatch(NameStart, parameter1, len(parameter2)){
			parameter1 = strings.Trim(prm1[2], ctlsp) // must be a Vmax
			parameter2 = strings.Trim(prm1[3], ctlsp) // must be a Hmax
			parameter3 = strings.Trim(prm1[4], ctlsp) // must be Starsmax
			parameter4 = strings.Trim(prm1[5], ctlsp) // must be a BHsmax
			parameter5 = strings.Trim(prm1[6], ctlsp) // must be a Archeronmax
			parameter6 = strings.Trim(prm1[7], ctlsp) // must be a Planetsmax
			parameter7 = strings.Trim(prm1[8], ctlsp) // must be a Port # in range
//fmt.Println("Post: got into start command!",parameter1,parameter2,parameter3,parameter4,parameter5,parameter6,parameter7)
			//
			// Validate the input
			//
			_, err = strconv.Atoi(parameter1)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			_, err = strconv.Atoi(parameter2)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			_, err = strconv.Atoi(parameter3)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			_, err = strconv.Atoi(parameter4)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			_, err = strconv.Atoi(parameter5)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			_, err = strconv.Atoi(parameter6)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			portnum, err = strconv.Atoi(parameter7)
			if err != nil {
				sendln(conn, invparm)
				return false
			}
			if (portnum < 1702) || (portnum > 1799) {
				sendln(conn, invparm)
				return false
			}

			//                 fmt.Println("starting with:",parameter1, parameter2, parameter3, parameter4, parameter5, parameter6, parameter7)
			//
			// Start another game
			//
//			a := "/usr/home/newmanh/go/src/decwars/decwars"
a := "/home/newmanh/go/src/decwars/decwars"
//			a := "/usr/home/newmanh/go/src/decwars/decwars " + parameter1 + " " + parameter2 +  " " + parameter3 +  " " +  parameter4 +  " " + parameter5 +  " " + parameter6 +  " " + parameter7 + " > ./" + parameter7 + "-" + "logfile.txt 2>&1 &"
//
//			a := "decwars > " + parameter7 + "-" + "logfile.txt 2>&1 &"
//			a := "decwars "
//			fmt.Println("here is the commandline: ",a)
			cmd := exec.Command(a, parameter1, parameter2, parameter3, parameter4, parameter5, parameter6, parameter7)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
			if err  != nil {
				fmt.Println("Got a error on fork:", err)
			}
//fmt.Println("Started game on port:", parameter7)			
			send(conn, "Started game on port:")
			sendln(conn, parameter7)

	} else {
		sendln(conn, invparm)
		return false
	}

}
		return false


}


//
// Bases command - with parms
//syntx := "BAses [ENemy | SUmmary | ALl | CLosest | RElative | ABsolute] location"
//
func processBases(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	var mySide int
	var myLocx int
	var myLocy int
	var myIOFmt int
	var myOutputLen int
	var parameter1 string
	var qry string
	var locx int
	var locy int
	var SeenbyEnemy int
	var Nme string
	// get my info
	objectsdb.QueryRow("select Side, Locx, Locy, IOFmt, OutputLen from objects where Nme=?", username).Scan(&mySide, &myLocx, &myLocy, &myIOFmt, &myOutputLen)

	prm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")
	//	prm1 := strings.SplitAfter(prm, " ")
	if len(prm) == 3 { // Must be a default address (choose rel/ab from user default)
		pival, err := strconv.Atoi(prm[1])

		if err == nil {
			locx = pival
			pival, err = strconv.Atoi(prm[2])
			if err == nil {
				locy = pival
				objectsdb.QueryRow("select IOFmt, Locx, Locy, Side from objects where Nme=?", username).Scan(&myIOFmt, &myLocx, &myLocy, &mySide)
			}
			if myIOFmt == IOFmtRel || myIOFmt == IOFmtBoth {
				if locx > MaxScanRng || locy > MaxScanRng {
					sendln(conn, OOR)
					return false
				}
				err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", myLocx+locx, myLocy+locy).Scan(&Nme, &SeenbyEnemy)
				if err == nil {
					istrue := SeenbyEnemy & mySide
					if istrue > 0 {
						qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(myLocx+locx) + " and Locy = " + strconv.Itoa(myLocy+locy) + " and Objtype = " + strconv.Itoa(TypeBase)
						doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
					} else {
						sendln(conn, OOR)
						return false
					}
				} else {
					sendln(conn, invparm)
					return false
				}
			} else { //either absolute
				if (Abs(myLocx)-Abs(locx) > MaxScanRng) || (Abs(myLocy)-Abs(locy) > MaxScanRng) {
					sendln(conn, OOR)
					return false
				}
				err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", locx, locy).Scan(&Nme, &SeenbyEnemy)
				if err == nil {
					istrue := SeenbyEnemy & mySide
					if istrue > 0 {
						qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(locx) + " and Locy = " + strconv.Itoa(locy) + " and Objtype = " + strconv.Itoa(TypeBase)
						doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
					} else {
						sendln(conn, OOR)
						return false
					}
				}
			}
		}
	} else {
		if len(prm) == 2 { // 1 parms - can be any keyword
			parameter1 = strings.Trim(prm[1], ctlsp)
			switch true {

			case testCommandMatch(NameEnemy, parameter1, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeBase)
				doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)

			case testCommandMatch(NameSummary, parameter1, lenOfString1):
				sendln(conn, " ")
				doSummaryBases(conn, objectsdb)

			case testCommandMatch(NameAll, parameter1, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + strconv.Itoa(TypeBase)
				doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)

			case testCommandMatch(NameClosest, parameter1, lenOfString1):
				clLocx, clLocy, rslt := findClosest(myLocx, myLocy, TypeBase, 0, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy) + " and Objtype = " + strconv.Itoa(TypeBase)
					doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
				} else {
					sendln(conn, "No closest base!")
					return true
				}
			}
		} else {
			if len(prm) == 4 { // 3 parms - can be any keyword
				parameter1 = strings.Trim(prm[1], ctlsp)
				switch true {
				case testCommandMatch(NameRelative, parameter1, lenOfString1):
					parameter2 := strings.Trim(prm[2], ctlsp)
					parameter3 := strings.Trim(prm[3], ctlsp)
					pival, err := strconv.Atoi(parameter2)
					if err == nil {
						locx := pival
						pival, err := strconv.Atoi(parameter3)
						if err == nil {
							locy := pival
							if locx > MaxScanRng || locy > MaxScanRng {
								sendln(conn, OOR)
								return false
							}
							qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(locx+myLocx) + " and Locy = " + strconv.Itoa(locy+myLocy) + " and Objtype = " + strconv.Itoa(TypeBase)
							doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
						}
						return false
					}

				case testCommandMatch(NameAbsolute, parameter1, lenOfString1):
					parameter2 := strings.Trim(prm[2], ctlsp)
					parameter3 := strings.Trim(prm[3], ctlsp)
					pival, err := strconv.Atoi(parameter2)
					if err == nil {
						locx := pival
						pival, err := strconv.Atoi(parameter3)
						if err == nil {
							locy := pival
							if (myLocx-locx) > MaxScanRng || (myLocy-locy) > MaxScanRng {
								sendln(conn, OOR)
								return false
							}
							qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(locx) + " and Locy = " + strconv.Itoa(locy) + " and Objtype = " + strconv.Itoa(TypeBase)
							doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
						}
						return false
					}

				default:
					sendln(conn, invparm)
					return true

				}
				sendln(conn, invparm)
				return false
			}
		}
		//					sendln(conn, invparm)
		return false
	}
	return false
}

//
// Do the default build command
//
func dodefBuild(invalidcommand bool, conn *Conn) bool {
	invalidcommand = processHelp("help build", invalidcommand, conn)
	return false
}

//
// Build command - with parms - requires adjacency to planet of your side.
// 5th build conditional on your side having < 10 bases
// bu computed
// bu x y
// bu rel/abs x y
// ccccccccccccccccccccccccc builds take 50 units of energy ?
//
func processBuild(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	var locx int
	var locy int
	var rslt bool

	var myLocx int
	var myLocy int
	var myShldUp int
	var myShipEnergy int
	var mySide int
	var myIOFmt int

	var whoLocx int
	var whoLocy int
	var whoShldUp int
	var whoShipEnergy int
	var whoNme string
	var whoType int
	var whoBuilds int
	var whoCount int
	var whoSide int
	var myTractorWho string

	var err error
	var parameter1 string
	var parameter2 string

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	// First parm
	parameter1 = strings.Trim(prm1[1], ctlsp)
	// Lookup me
	objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy, IOFmt, Side, TractorWho from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myShipEnergy, &myIOFmt, &mySide, &myTractorWho)

	if len(prm1) == 2 { // 1 parm (bu computed)
		if testCommandMatch(NameComputed, parameter1, lenOfString1) {
			locx, locy, rslt = docomputed(username, NamePlanets, objectsdb)
			if rslt {
				parameter1 = strconv.Itoa(locx)
				parameter2 = strconv.Itoa(locy)
				myIOFmt = IOFmtAbs
			} else {
				sendln(conn, nAdjPlt)
				return false
			}
		} else {
			sendln(conn, invparm)
			return false
		}
	}
	if len(prm1) == 3 { // 2 parm - bu x y
		parameter2 = strings.Trim(prm1[2], ctlsp)

		// Lookup user
		locx, _ = strconv.Atoi(parameter1)
		locy, _ = strconv.Atoi(parameter2)
		if myIOFmt == IOFmtAbs {
			err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
		} else {
			err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
		}
		if whoType != TypePlanet {
			sendln(conn, nAdjPlt)
			return false
		}
		if err != nil { // no such obj
			sendln(conn, invparm)
			return false
		} else {
			locx = whoLocx
			locy = whoLocy
		}
	}
	if len(prm1) == 4 { // 3 parm - bu r/a x y

		locx, _ = strconv.Atoi(strings.Trim(prm1[2], ctlsp))
		locy, _ = strconv.Atoi(strings.Trim(prm1[3], ctlsp))

		switch true {
		case testCommandMatch(NameRelative, parameter1, lenOfString1):
			myIOFmt = IOFmtRel
			err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy, Objtype from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
			break

		case testCommandMatch(NameAbsolute, parameter1, lenOfString1):
			myIOFmt = IOFmtAbs
			err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy, Objtype from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
			break

		default:
			sendln(conn, nAdjPlt)
			return false

		}

		if whoType != TypePlanet {
			sendln(conn, nAdjPlt)
			return false
		}

		locx = whoLocx
		locy = whoLocy
	}

	// so if we get here locx, locy = valid obj, is it adjacent and on our side?

	if (Abs(myLocx-locx) <= 1) && (Abs(myLocy-locy) <= 1) {
		objectsdb.QueryRow("select Nme, Builds, Side from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoBuilds, &whoSide)
		// if the number of builds > 4 convert to base ... if less than 10 bases
		if whoSide != mySide {
			sendln(conn, nEnRefuse)
			return false
		}
		if whoBuilds > 3 {
			objectsdb.QueryRow("select count(Nme) from objects WHERE ObjType = ? and Side = ?", TypeBase, mySide).Scan(&whoCount)
			if whoCount <= Basesmax {
				objectsdb.Exec("UPDATE objects set ObjType = ?, Builds = ?, ShldUp = ? where Locx = ? and Locy = ?", TypeBase, 0, On, locx, locy)
				//				send(conn, username)
				//				sendln(conn, BPB)
				notify(username, conn, myLocx, myLocx, msgBuildBase, username, myLocx, myLocx, objectsdb, 0, 0)
				//
				// if you have a tractor beam to the base, turn off the tractor and notify bases shields raised, tractor broken
				//
				if myTractorWho == whoNme {
					endTractor(username, conn, objectsdb)
					sendln(conn, "Base raised shields, captain, tractor beam broken!")
				}
				return false
			} else {
				sendln(conn, ABF)
				return false
			}
		} else {
			objectsdb.Exec("UPDATE objects set Builds = Builds + 1 where Locx = ? and Locy = ?", locx, locy)
			//
			// Tell user he built it
			if whoBuilds == 0 {
				sendln(conn, "1 build")
			} else {
				send(conn, strconv.Itoa(whoBuilds+1))
				sendln(conn, " builds")
			}
		}
		// build in a delay
		doGameDelay(username, whoBuilds+1, objectsdb)
	} else {
		sendln(conn, nAdjPlt)
		return false
	}
	return false
}

//
// Do the default Capture command
//
func dodefCapture(invalidcommand bool, conn *Conn) bool {
	invalidcommand = processHelp("help capture", invalidcommand, conn)
	return false
}

//
// Capture command with parms - require adjacency
// ca x y
// ca computed side
// ca re/ab x y
//
func processCapture(username string, comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var myLocx int
	var myLocy int
	var myShldUp int
	var myShipEnergy int
	var myIOFmt int
	var mySide int
	var parameter1 string
	var parameter2 string
	var locx int
	var locy int
	var whoNme string
	var whoLocx int
	var whoLocy int
	var whoShldUp int
	var whoShipEnergy int
	var whoType int
	var err error
	var err1 error
	var rslt bool
	var bld int

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	// First parm
	parameter1 = strings.Trim(prm1[1], ctlsp)
	// Lookup me
	objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy, IOFmt, Side from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myShipEnergy, &myIOFmt, &mySide)

	//fmt.Println("going into test for 4 (ab or rel)", len(prm1))

	//
	// If parms = 4 (cap rel 4 4 or ca ab 4 4)
	//
	if len(prm1) == 4 {
		parameter2 = strings.Trim(prm1[2], ctlsp)
		parameter3 := strings.Trim(prm1[3], ctlsp)
		// Convert to numeric
		locx, err = strconv.Atoi(parameter2)
		locy, err1 = strconv.Atoi(parameter3)

		//fmt.Println(" 4 parms, parmater1=", parameter1, " parameter2=", parameter2, " parameter3=", parameter3, "locx=", locx, " locy:",locy)

		if testCommandMatch(NameAbsolute, parameter1, lenOfString1) {
			//fmt.Println("got into absolute, locx:", locx, "locy:",locy)
			myIOFmt = IOFmtAbs
			err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype, Builds from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType, &bld)
		}
		if testCommandMatch(NameRelative, parameter1, lenOfString1) {
			myIOFmt = IOFmtRel
			err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype, Builds from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType, &bld)
			//fmt.Println("got into rel, locx:", locx, "locy:",locy, "mlocx:", myLocx+locx, "mlocy:",myLocy+locy, " mylocy:",myLocy, " mylocx:",myLocx)
		}
	}

	//
	// Ca com pl = 3 parms all alpha
	//
	//fmt.Println("entering processcapture prm1:", prm1, " parameter1:", parameter1," len:",len(prm1))

	if len(prm1) == 3 { // 2 parm - ca com pl OR name
		parameter2 = strings.Trim(prm1[2], ctlsp)

		// Convert to numeric
		locx, err = strconv.Atoi(parameter1)
		locy, err1 = strconv.Atoi(parameter2)

		// if first parm is alpha, then 2nd must be alpha (for ca comp planet or name) (err and err1 !=nil)
		//fmt.Println("err; ",err, " err1: ",err1, " parameter1:", parameter1)
		if err != nil && err1 != nil {
			//fmt.Println("got two strings")
			if testCommandMatch(NameComputed, parameter1, lenOfString1) {
				//fmt.Println("got computed")
				if testCommandMatch(NamePlanets, parameter2, len(parameter2)) {
					locx, locy, rslt = docomputed(username, NamePlanets, objectsdb)
					//fmt.Println("got planets,  x-",locx, " y--",locy, " rslt:", rslt)
					parameter1 = strconv.Itoa(locx)
					parameter2 = strconv.Itoa(locy)
					if rslt {
						// get info for computed
						err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype, Builds from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType, &bld)
					}
				} else { //specified name of planet
					locx, locy, rslt = docomputed(username, parameter2, objectsdb)
					//fmt.Println("got name of planet,  x-",locx, " y--",locy, " rslt:", rslt)
					if rslt {
						parameter1 = strconv.Itoa(locx)
						parameter2 = strconv.Itoa(locy)
						// get info for computed
						err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype, Builds from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType, &bld)
					} else {
						sendln(conn, invparm)
					}
				}
			}
		} else {

			//fmt.Println("going into default test ", len(prm1))
			//
			// parms must be 3 or error
			//
			if len(prm1) == 3 {
				//fmt.Println("testing ioformat")
				if myIOFmt == IOFmtAbs {
					err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype, Builds from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType, &bld)
				} else {
					err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype, Builds from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType, &bld)
				}
				if err != nil { // no such obj
					sendln(conn, noTarg)
					return false
				} else {
					locx = whoLocx
					locy = whoLocy
				}
			}

		}
	}
	//fmt.Println(" whotype:", whoType, " whoLocx:",whoLocx," whoLocy:",whoLocy)
	if whoType != TypePlanet {
		sendln(conn, nAPlt)
		return false
	}
	if (Abs(myLocx-whoLocx) > 1) || (Abs(myLocy-whoLocy) > 1) {
		sendln(conn, nAdjPlt)
		return false
	}

	//fmt.Println("doing capture")
	//
	// Right here do the capture - simple - change side from neutral to your side - possibility to fail!
	//
	if rand.Float32() <= .85 {
		objectsdb.Exec("UPDATE objects set Side = ? WHERE Nme = ?", mySide, whoNme)
		sendln(conn, plCaptured)
	} else {
		sendln(conn, nRefuseCapture)
	}

	//
	// Delay here
	//
	doGameDelay(username, 1+bld, objectsdb)
	return false
}

//
// Print damage header
//
func doDamagesHeader(conn *Conn, username string) {
	send(conn, "Damage report for ")
	sendln(conn, username)
	return
}

//
// Do the default damages command
//
func dodefDamages(conn *Conn, username string, objectsdb *sql.DB) bool {
	doDamagesHeader(conn, username)
	doWarpDam(conn, username, objectsdb)
	doImpulseDam(conn, username, objectsdb)
	doTorpedoesDam(conn, username, objectsdb)
	doPhasersDam(conn, username, objectsdb)
	doShieldDam(conn, username, objectsdb)
	doComputerDam(conn, username, objectsdb)
	doLifeSupportDam(conn, username, objectsdb)
	doRadioDam(conn, username, objectsdb)
	doTractorDam(conn, username, objectsdb)
	doShipDam(conn, username, objectsdb)
	return false
}
func doWarpDam(conn *Conn, username string, objectsdb *sql.DB) {
	var WarpEngDam int

	_ = objectsdb.QueryRow("select WarpEngDam from objects WHERE Nme = ?", username).Scan(&WarpEngDam)

	send(conn, "Warp engine damage:\t\t")
	sendln(conn, strconv.Itoa(WarpEngDam))
	return
}

func doImpulseDam(conn *Conn, username string, objectsdb *sql.DB) {
	var ImpEngDam int

	_ = objectsdb.QueryRow("select ImpEngDam from objects WHERE Nme = ?", username).Scan(&ImpEngDam)

	send(conn, "Impulse engine damage:\t\t")
	sendln(conn, strconv.Itoa(ImpEngDam))
	return
}

func doTorpedoesDam(conn *Conn, username string, objectsdb *sql.DB) {
	var PhoTorDam int

	_ = objectsdb.QueryRow("select PhoTorDam from objects WHERE Nme = ?", username).Scan(&PhoTorDam)

	send(conn, "Photon torpedoes damage:\t")
	sendln(conn, strconv.Itoa(PhoTorDam))
	return
}

func doPhasersDam(conn *Conn, username string, objectsdb *sql.DB) {
	var PhasDam int

	_ = objectsdb.QueryRow("select PhasDam from objects WHERE Nme = ?", username).Scan(&PhasDam)

	send(conn, "Phaser bank damage:\t\t")
	sendln(conn, strconv.Itoa(PhasDam))
	return
}

func doShieldDam(conn *Conn, username string, objectsdb *sql.DB) {
	var ShldDam int

	_ = objectsdb.QueryRow("select ShldDam from objects WHERE Nme = ?", username).Scan(&ShldDam)

	send(conn, "Shield damage:\t\t\t")
	sendln(conn, strconv.Itoa(ShldDam))
	return
}

func doComputerDam(conn *Conn, username string, objectsdb *sql.DB) {
	var CmpDam int

	_ = objectsdb.QueryRow("select CmpDam from objects WHERE Nme = ?", username).Scan(&CmpDam)

	send(conn, "Computer damage:\t\t")
	sendln(conn, strconv.Itoa(CmpDam))
	return
}

func doLifeSupportDam(conn *Conn, username string, objectsdb *sql.DB) {
	var LifeSupDam int

	_ = objectsdb.QueryRow("select LifeSupDam from objects WHERE Nme = ?", username).Scan(&LifeSupDam)

	send(conn, "Life support damage:\t\t")
	sendln(conn, strconv.Itoa(LifeSupDam))
	return
}

func doRadioDam(conn *Conn, username string, objectsdb *sql.DB) {
	var RadioDam int

	_ = objectsdb.QueryRow("select RadioDam from objects WHERE Nme = ?", username).Scan(&RadioDam)

	send(conn, "Radio damage:\t\t\t")
	sendln(conn, strconv.Itoa(RadioDam))
	return
}

func doTractorDam(conn *Conn, username string, objectsdb *sql.DB) {
	var TractorDam int

	_ = objectsdb.QueryRow("select TractorDam from objects WHERE Nme = ?", username).Scan(&TractorDam)

	send(conn, "Tractor beam damage:\t\t")
	sendln(conn, strconv.Itoa(TractorDam))
	return
}

func doShipDam(conn *Conn, username string, objectsdb *sql.DB) {
	var ShipDam int

	_ = objectsdb.QueryRow("select ShipDam from objects WHERE Nme = ?", username).Scan(&ShipDam)

	send(conn, "Total ship damage:\t\t")
	sendln(conn, strconv.Itoa(ShipDam))
	return
}

//
// Damages command - DAmages [device | username | RElative | ABsolute] location]
//
func processDamages(comnd string, username string, conn *Conn, objectsdb *sql.DB) bool {
	var Nme string
	var SeenbyEnemy int
	var mySide int
	var myIOFmt int
	var myLocx int
	var myLocy int
	var locx int
	var locy int

	prm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")

	if len(prm) == 2 {
		//		doDamagesHeader(conn, username)
		switch true {
		case testCommandMatch("warp", prm[1], lenOfString1):
			doWarpDam(conn, username, objectsdb)

		case testCommandMatch("impulse", prm[1], lenOfString1):
			doImpulseDam(conn, username, objectsdb)

		case testCommandMatch("photon torpedoes", prm[1], lenOfString3):
			doTorpedoesDam(conn, username, objectsdb)

		case testCommandMatch(NamePhas, prm[1], lenOfString3):
			doPhasersDam(conn, username, objectsdb)

		case testCommandMatch("shields", prm[1], lenOfString3):
			doShieldDam(conn, username, objectsdb)

		case testCommandMatch("computer", prm[1], lenOfString1):
			doComputerDam(conn, username, objectsdb)

		case testCommandMatch("life support", prm[1], lenOfString1):
			doLifeSupportDam(conn, username, objectsdb)

		case testCommandMatch("radio", prm[1], lenOfString1):
			doRadioDam(conn, username, objectsdb)

		case testCommandMatch("tractor", prm[1], lenOfString2):
			doTractorDam(conn, username, objectsdb)

		case testCommandMatch("ship", prm[1], lenOfString3):
			doShipDam(conn, username, objectsdb)

		default: //  Must be username
			objectsdb.QueryRow("select Side from objects where Nme=?", username).Scan(&mySide)
			err := objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Nme=?", prm[1]).Scan(&Nme, &SeenbyEnemy)
			if err == nil {
				istrue := SeenbyEnemy & mySide
				if istrue > 0 { //must have been seen to see damages!  <= Works!
					dodefDamages(conn, prm[1], objectsdb)
				} else {
					sendln(conn, "Out of range")
					return false
				}
			} else {
				sendln(conn, "No such user!")
				return false
			}
		}
	} else {
		if len(prm) == 3 { // Must be a default address (choose rel/ab from user default)
			pival, err := strconv.Atoi(prm[1])
			if err == nil {
				locx = pival
				pival, err = strconv.Atoi(prm[2])
				if err == nil {
					locy = pival
					objectsdb.QueryRow("select IOFmt, Locx, Locy, Side from objects where Nme=?", username).Scan(&myIOFmt, &myLocx, &myLocy, &mySide)
				}
				if myIOFmt == IOFmtRel || myIOFmt == IOFmtBoth {
					if locx > MaxScanRng || locy > MaxScanRng {
						sendln(conn, OOR)
						return false
					}
					err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", myLocx+locx, myLocy+locy).Scan(&Nme, &SeenbyEnemy)
					if err == nil {
						istrue := SeenbyEnemy & mySide
						if istrue > 0 {
							dodefDamages(conn, Nme, objectsdb)
						} else {
							sendln(conn, OOR)
							return false
						}
					} else {
						sendln(conn, invparm)
						return false
					}
				} else { // absolute
					if (Abs(myLocx)-Abs(locx) > MaxScanRng) || (Abs(myLocy)-Abs(locy) > MaxScanRng) {
						sendln(conn, OOR)
						return false
					}
					err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", locx, locy).Scan(&Nme, &SeenbyEnemy)
					if err == nil {
						istrue := SeenbyEnemy & mySide
						if istrue > 0 {
							dodefDamages(conn, Nme, objectsdb)
						} else {
							sendln(conn, OOR)
							return false
						}
					}
				}
			}
		} else {
			if len(prm) == 4 {
				if testCommandMatch(NameRelative, prm[1], lenOfString1) {
					pival, err := strconv.Atoi(prm[2])
					if err == nil {
						locx := pival
						pival, err := strconv.Atoi(prm[3])
						if err == nil {
							locy := pival
							if locx > MaxScanRng || locy > MaxScanRng {
								sendln(conn, OOR)
								return false
							}
							objectsdb.QueryRow("select Locx, Locy, Side from objects where Nme=?", username).Scan(&myLocx, &myLocy, &mySide)
							err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", myLocx+locx, myLocy+locy).Scan(&Nme, &SeenbyEnemy)
							if err == nil {
								istrue := SeenbyEnemy & mySide
								if istrue > 0 {
									dodefDamages(conn, Nme, objectsdb)
								} else {
									sendln(conn, OOR)
									return false
								}
							}
						} else {
							sendln(conn, invparm)
							return false
						}
					} else {
						sendln(conn, invparm)
						return false
					}
				} else {
					if testCommandMatch("absolute", prm[1], lenOfString1) {
						pival, err := strconv.Atoi(prm[2])
						if err == nil {
							locx := pival
							pival, err := strconv.Atoi(prm[3])
							if err == nil {
								locy := pival
								objectsdb.QueryRow("select Side from objects where Nme=?", username).Scan(&mySide)

								err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", locx, locy).Scan(&Nme, &SeenbyEnemy)
								if err == nil {
									istrue := SeenbyEnemy & mySide
									if istrue > 0 {
										dodefDamages(conn, Nme, objectsdb)
									} else {
										sendln(conn, OOR)
										return false
									}
								}
							} else {
								sendln(conn, invparm)
								return false
							}
						} else {
							sendln(conn, invparm)
							return false
						}
					} else { // use default iofmt setting must be numeric
						sendln(conn, invparm)
						return false
					}
				}
			} else {
				sendln(conn, invparm)
				return false
			}
		}
	}
	return false
}

//
// Do the default disconnect command
//
func dodefDisconnect(username string, invalidcommand bool, conn *Conn) bool {
	sendln(conn, "Thanks for playing...")
	conn.Close()
	if username != "" {
		// mu.Lock()
		dbDelete(strings.Trim(username, ctlsp))
		// mu.Unlock()
		runtime.Gosched()
	}
	return false
}

//
// Function to end a dock.
//
func endDock(username string, conn *Conn, objectsdb *sql.DB) {
	var myDockFlag string

	//see if we already are tractoring tractored & end
	objectsdb.QueryRow("select DockFlag from objects WHERE Nme = ?", username).Scan(&myDockFlag)

	//Assume we are docked if we have this function called
	objectsdb.Exec("UPDATE objects set DockFlag = ? WHERE Nme = ?", "", username)
	sendln(conn, endDockMsg)
	return
}

//
// Do the default dock command - no parms - same as dock closest port
//
func dodefDock(username string, comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var MyDockFlag string
	objectsdb.QueryRow("select DockFlag from objects WHERE Nme = ?", username).Scan(&MyDockFlag)
	if MyDockFlag == "" {
		processDock(username, "dock closest planet", invalidcommand, conn, objectsdb, true)
		objectsdb.QueryRow("select DockFlag from objects WHERE Nme = ?", username).Scan(&MyDockFlag)
		if MyDockFlag == "" {
			processDock(username, "dock closest base", invalidcommand, conn, objectsdb, true)
			objectsdb.QueryRow("select DockFlag from objects WHERE Nme = ?", username).Scan(&MyDockFlag)
			if MyDockFlag == "" {
				sendln(conn, nAdjPort)
			}
		}
	} else {
		doGameDelay(username, 1, objectsdb)
		sendln(conn, successDock)
	}
	return true
}

//
// Dock command
//
func processDock(username string, comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, quiteAdjacency bool) bool {
	var Nme string
	var Locx int
	var Locy int
	var MySide int
	var MyDockFlag string
	var IOFmt int
	var WarpEngDam int
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var rslt bool
	var aNme string
	var aObjtype int
	var aSide int
	var p2 int
	var p3 int
	var newlocx int
	var newlocy int

	//   Move [Absolute|Relative|Computed] <vpos> <hpos> | Planet | Base
	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	if len(prm1) > 4 {
		sendln(conn, "Too many parameters")
		return true
	}
	if len(prm1) < 3 {
		sendln(conn, "Too few parameters")
		return true
	}

	//
	// Get the user's data
	//
	objectsdb.QueryRow("select Nme, Locx, Locy, IOFmt, WarpEngDam, Side, DockFlag from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &IOFmt, &WarpEngDam, &MySide, &MyDockFlag)

	//
	//if len(prm1) = 3 then both parms must be int (x,y) OR both are string (closest Planet|Base)
	//
	if len(prm1) == 3 {
		parameter1 = strings.Trim(prm1[1], ctlsp)
		parameter2 = strings.Trim(prm1[2], ctlsp)
		if _, err := strconv.Atoi(parameter1); err == nil {
			if _, err := strconv.Atoi(parameter2); err != nil {
				sendln(conn, invparm)
				return true
			}
		}
		if _, err := strconv.Atoi(parameter1); err != nil {
			if _, err := strconv.Atoi(parameter2); err == nil {
				sendln(conn, invparm)
				return true
			}
		}
		if IOFmt == IOFmtAbs {
			p2, _ = strconv.Atoi(parameter1)
			p3, _ = strconv.Atoi(parameter2)

			parameter2 = strconv.Itoa(p2 - Locx)
			parameter3 = strconv.Itoa(p3 - Locy)
		} else {
			parameter3 = parameter2
			parameter2 = parameter1
		}
	}

	//
	// If len(prm1) = 4 then 1st parm must be absolute and 2nd & 3rd must be numbers
	//
	if len(prm1) == 4 {
		parameter1 = strings.Trim(prm1[1], ctlsp)
		parameter2 = strings.Trim(prm1[2], ctlsp)
		parameter3 = strings.Trim(prm1[3], ctlsp)

		if _, err := strconv.Atoi(parameter1); err == nil {
			sendln(conn, invparm)
			return true
		}
		if _, err := strconv.Atoi(parameter2); err != nil {
			sendln(conn, "invparm")
			return true
		}
		if _, err := strconv.Atoi(parameter3); err != nil {
			sendln(conn, invparm)
			return true
		}
	}

	//
	// Was the request for absolute - if so, convert parms to relative
	//
	//		pmx := parameter2
	//		pmy := parameter3
	if testCommandMatch("absolute", parameter1, lenOfString1) == true {
		// force users format to absolute
		IOFmt = IOFmtAbs
		p2, _ = strconv.Atoi(parameter2)
		p3, _ = strconv.Atoi(parameter3)
		parameter2 = strconv.Itoa(p2 - Locx)
		parameter3 = strconv.Itoa(p3 - Locy)
	} else {
		if testCommandMatch(NameRelative, parameter1, lenOfString1) == true {
			// force users format to relative
			IOFmt = IOFmtRel
		} else { //closest- dock closest matching object - compute the loc and force it to absolute!
			if testCommandMatch(NameClosest, parameter1, lenOfString1) == true {
				// need to use computer here
				newlocx, newlocy, rslt = docomputed(username, parameter3, objectsdb)
				if rslt {
					parameter2 = strconv.Itoa(newlocx - Locx)
					parameter3 = strconv.Itoa(newlocy - Locy)
					IOFmt = IOFmtAbs
				} else {
					sendln(conn, "No port within range, Captain!")
					return false
				}
			}
		}
	}

	//
	// Test for adjacency
	//
	p2, _ = strconv.Atoi(parameter2)
	p3, _ = strconv.Atoi(parameter3)
	//
	if (Abs(p2) > 1) || (Abs(p3) > 1) {
		if quiteAdjacency == false {
			sendln(conn, nAdjPort)
		}
		return false
	}

	//
	// Is there an object, and is it a port?
	//
	err := objectsdb.QueryRow("select Nme, Objtype, Side from objects WHERE Locx = ? and Locy = ?", Locx+p2, Locy+p3).Scan(&aNme, &aObjtype, &aSide)
	if err != nil {
		if quiteAdjacency == false {
			sendln(conn, nAdjPort)
		}
		return false
	}
	if ((aObjtype != TypePlanet) && (aObjtype != TypeBase)) || (aSide != MySide) {
		if quiteAdjacency == false {
			sendln(conn, nAdjPort)
		}
		return false
	}
	//
	// Now we can do the dock
	//
	objectsdb.Exec("UPDATE objects set DockFlag = ?  WHERE Nme = ?", aNme, username)
	doGameDelay(username, 1, objectsdb)
	sendln(conn, successDock)
	return false
}

//
// Energy command -with parms:
// en ship amt
// en loc loc amt
// en [rel/abs] loc loc amt
//
func processEnergy(username string, comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var myLocx int
	var myLocy int
	var myShldUp int
	var myShipEnergy int
	var mySide int
	var myIOFmt int
	var whoLocx int
	var whoLocy int
	var whoShldUp int
	var whoShipEnergy int
	var whoNme string
	var whoObjType int
	var whoSide int
	var amtToTransfer int
	var prm2 int
	var prm3 int
	var prm4 int
	var err error
	var locx int
	var locy int

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	amtToTransfer = 0
	// First parm
	parameter1 := strings.Trim(prm1[1], ctlsp)

	if len(prm1) == 3 { // 2 parm (name amount)
		// Lookup me
		objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy, Side from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myShipEnergy, &mySide)
		// Lookup user
		whoNme = parameter1
		err = objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy from objects WHERE Nme = ?", whoNme).Scan(&whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy)

		if err == nil { // parm must be shipname, parm 2 must be amount
			parameter2 := strings.Trim(prm1[2], ctlsp)
			if prm2, err = strconv.Atoi(parameter2); err != nil {
				sendln(conn, invparm)
				return true
			}
			//
			// Object must be adjacent and shields down for both AND we are not tractoring someone, and they are not being tractored!
			//
			if (Abs(myLocx-whoLocx) <= 1) && (Abs(myLocy-whoLocy) <= 1) {
				if (myShldUp == Off) && (whoShldUp == Off) {
					whoMaxNeededEnergy := InitEnergy - whoShipEnergy
					myMaxNeededEnergy := InitEnergy - myShipEnergy
					// is this a pull or push
					if prm2 <= 0 { //pull
						if Abs(prm2) >= myMaxNeededEnergy {
							amtToTransfer = -myMaxNeededEnergy
						} else {
							amtToTransfer = prm2
						}
					} else { //push
						if prm2 >= whoMaxNeededEnergy {
							amtToTransfer = whoMaxNeededEnergy
						} else {
							amtToTransfer = prm2
						}
					}
				} else {
					sendln(conn, "Shields must be down for both ships, Captain")
					return false
				}
			} else {
				send(conn, "The ")
				send(conn, parameter2)
				sendln(conn, " is not adjacent")
				return false

			}
		}
	} else {
		if len(prm1) == 4 { // 3 parm (mo x y amt)
			parameter2 := strings.Trim(prm1[2], ctlsp)
			parameter3 := strings.Trim(prm1[3], ctlsp)
			prm3, _ := strconv.Atoi(parameter3)

			// Lookup me
			objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy, IOFmt, Side from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myShipEnergy, &myIOFmt, &mySide)
			// Lookup user
			locx, _ := strconv.Atoi(parameter1)
			locy, _ := strconv.Atoi(parameter2)
			if myIOFmt == IOFmtAbs {
				err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy)
			} else {
				err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy)
			}

			if err == nil { // parm 4 must be amount
				if prm2, err = strconv.Atoi(parameter3); err != nil {
					sendln(conn, invparm)
					return true
				}
				//
				// Object must be adjacent and shields down for both AND we are not tractoring someone, and they are not being tractored!
				//
				if (Abs(myLocx-whoLocx) <= 1) && (Abs(myLocy-whoLocy) <= 1) {
					if (myShldUp == 0) && (whoShldUp == 0) {
						whoMaxNeededEnergy := InitEnergy - whoShipEnergy
						myMaxNeededEnergy := InitEnergy - myShipEnergy
						// is this a pull or push
						if prm2 <= 0 { //pull
							if Abs(prm3) >= myMaxNeededEnergy {
								amtToTransfer = -myMaxNeededEnergy
							} else {
								amtToTransfer = prm3
							}
						} else { //push
							if prm3 >= whoMaxNeededEnergy {
								amtToTransfer = whoMaxNeededEnergy
							} else {
								amtToTransfer = prm3
							}
						}
					} else {
						sendln(conn, "Shields must be down for both ships, Captain")
						return false
					}
				} else {
					send(conn, "The ")
					send(conn, whoNme)
					sendln(conn, " is not adjacent")
					return false

				}
			}
		} else {
			if len(prm1) == 5 { // 4 parm  (en r/a x y amt)
				parameter2 := strings.Trim(prm1[2], ctlsp)
				prm2, _ = strconv.Atoi(parameter2)
				parameter3 := strings.Trim(prm1[3], ctlsp)
				prm3, _ = strconv.Atoi(parameter3)
				parameter4 := strings.Trim(prm1[4], ctlsp)
				// Lookup me
				objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy, IOFmt, Side from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myShipEnergy, &myIOFmt, &mySide)
				// Lookup user
				locx = prm2
				locy = prm3
				if testCommandMatch(NameAbsolute, parameter1, lenOfString1) {
					err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy)
				} else {
					err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy)
				}
				if err == nil { // parm 4 must be amount
					if prm4, err = strconv.Atoi(parameter4); err != nil {
						sendln(conn, invparm)
						return true
					}
					//
					// Object must be adjacent and shields down for both AND we are not tractoring someone, and they are not being tractored!
					//
					if (Abs(myLocx-whoLocx) <= 1) && (Abs(myLocy-whoLocy) <= 1) {
						if (myShldUp == 0) && (whoShldUp == 0) {
							whoMaxNeededEnergy := InitEnergy - whoShipEnergy
							myMaxNeededEnergy := InitEnergy - myShipEnergy
							// is this a pull or push
							if prm4 <= 0 { //pull
								if Abs(prm4) >= myMaxNeededEnergy {
									amtToTransfer = -myMaxNeededEnergy
								} else {
									amtToTransfer = prm4
								}
							} else { //push
								if prm4 >= whoMaxNeededEnergy {
									amtToTransfer = whoMaxNeededEnergy
								} else {
									amtToTransfer = prm4
								}
							}
						} else {
							sendln(conn, "Shields must be down for both ships, Captain")
							return false
						}
					} else {
						send(conn, "The ")
						send(conn, whoNme)
						sendln(conn, " is not adjacent")
						return false

					}
				}
			}
		}
	}

	if username == whoNme {
		sendln(conn, "Transfer energy to US!?!")
		return false
	}

	objectsdb.QueryRow("select Objtype, Side from objects WHERE Nme = ?", whoNme).Scan(&whoObjType, &whoSide)
	if whoObjType != TypeShip {
		sendln(conn, "Not a ship, captain!")
		return false
	}

	if whoSide != mySide {
		sendln(conn, "Can not transfer energy to enemy ship, captain!")
		return false
	}

	degradedEnergy := int(.9 * float64(amtToTransfer))

	_, err = objectsdb.Exec("UPDATE objects set ShipEnergy = ? WHERE Nme = ?", whoShipEnergy+degradedEnergy, whoNme)

	//10% of the energy transferred will be lost due to broadcast dissipation.
	if amtToTransfer != 0 {
		_, err = objectsdb.Exec("UPDATE objects set ShipEnergy = ? WHERE Nme = ?", myShipEnergy-amtToTransfer, username)
		_, err = objectsdb.Exec("UPDATE objects set ShipEnergy = ? WHERE Nme = ?", whoShipEnergy+degradedEnergy, whoNme)
	} else {
		sendln(conn, "Cannot transfer 0!")
	}
	return false
}

//
// Do the default energy command - without parms
//
func dodefEnergy(username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var ShipEnergy int
	//username is currently mine!!
	_ = objectsdb.QueryRow("select ShipEnergy from objects WHERE Nme = ?", username).Scan(&ShipEnergy)
	//
	doStatusEnergy(conn, ShipEnergy)
	return false
}

//
// Do the default gripe command
//
func dodefGripe(invalidcommand bool, conn *Conn) bool {
	// Open file for reading
	file, err := os.Open("/home/newmanh/go/src/decwars/gripe.txt")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	s := string(data[:])
	sendln(conn, s)
	return false
}

//
// Gripe command
//
func processGripe(comnd string, invalidcommand bool, conn *Conn, usrnme string, netaddr string) bool {
	msg := strings.Split(strings.TrimSpace(strings.Trim(comnd, ctlonly)), " ")
	if len(msg) > 1 {
		// Open the file
		f, err := os.OpenFile("/home/newmanh/go/src/decwars/gripe.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		} //
		n, err := f.WriteString(time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"))
		if err != nil {
			log.Fatal(err)
		}
		n, err = f.WriteString(" ")
		if err != nil {
			log.Fatal(err)
		}
		n, err = f.WriteString(usrnme)
		if err != nil {
			log.Fatal(err)
		}
		n, err = f.WriteString(", ")
		if err != nil {
			log.Fatal(err)
		}
		n, err = f.WriteString(netaddr)
		if err != nil {
			log.Fatal(err)
		}
		n, err = f.WriteString(": ")
		if err != nil {
			log.Fatal(err)
		}
		n, err = f.WriteString(strings.Join(msg[1:], " "))
		if err != nil {
			log.Fatal(err)
		}
		n, err = f.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}
		n = n
		f.Close()
		sendln(conn, "Gripe accepted.\n")
	}
	return false
}

//
// Help command
//
func processHelp(comnd string, invalidcommand bool, conn *Conn) bool {
	sndto := strings.Split(strings.TrimSpace(strings.Trim(comnd, ctlonly)), " ")
	// just asked for help
	if len(sndto) == 2 {
		if strings.ToLower(sndto[1]) == "?" || strings.ToLower(sndto[1]) == "*" {
			invalidcommand = dodefHelp(invalidcommand, conn)
		} else {
			switch true {
			case testCommandMatch("bases", sndto[1], lenOfString2):
				nme := "BAses"
				syntx := "BAses [ENemy | SUmmary | ALl | CLosest | RElative | ABsolute] <Locy Locx>]"
				fctn := "List information on friendly and  known  enemy bases."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "BAses\t\t\t\t(List all known bases with summary)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BAses ENemy\t\t\t(List all known enemy bases)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BAses SUmmary\t\t\t(List a summary of all known bases)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BAses ALl\t\t\t(List all known bases)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BAses CLosest\t\t\t(List closest knwon base)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BAses RElative -3 5\t\t(List the base relative to your ship -3 5)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BAses ABsolute 37 53\t\t(List the base at the location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("build", sndto[1], lenOfString2):
				nme := "BUild"
				syntx := "BUild [<RElative | ABsolute | Computed>] <Locy Locx>"
				fctn := "Develop  installations  on  a  planet,   and eventually build it into a base.  The planet must first be captured."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "BUild RElative 1 -1\t\t(Build the planet 1 -1 from you)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BUild ABsolute 37 53\t\t(Build the planet at the location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "BUild Computed p47\t\t(Build the planet #47)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("capture", sndto[1], lenOfString1):
				nme := "Capture"
				syntx := "Capture [<Computed | RElative | Absolute>] <Locy Locx>"
				fctn := "Win a neutral or enemy planet over  to  your side."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "CApture COmputed p47\t\t(Capture the planet named p47)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "CApture Relative 1 -1\t\t(Capture the planet relative to your ship 1 -1)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "CApture ABsolute 37 53\t\t(Capture the planet at the location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("damages", sndto[1], lenOfString2):
				nme := "DAmages"
				syntx := "DAmages [device | username | RElative | ABsolute] location]"
				fctn := "List  damaged  devices  and  their  current status."
				optns := "device | username | RElative | ABsolute"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "DAmages WArp\t\t(List damages to warp engines.)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages SHIelds\t\t(List damages to shields)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages  PHAser\t\t(List damages to Phaser banks)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages COmputer\t(List damages to computer)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages LIfe\t\t(List damages to life support)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages RAdio\t\t(List damages to radio)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages SHip\t\t(List total damages to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages A3\t\t(List damages to Archeron #3)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages Relative 1 1\t(List damages to object located relative 1 1 to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DAmages Absolute 34 37\t(List damages to object located at 34 37)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("dock", sndto[1], lenOfString2):
				nme := "DOck"
				syntx := "DOck [Absolute|Relative|CLosest] <RElative | ABsolute> | Planet | Base"
				fctn := "Dock at  an  adjacent  base  or  planet.   This increases  your  energy,  replenishes  your  torpedoes, repairs your ship  a  little,  and  reduces  your  ship damage."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "DOck\t\t\t\t(Dock to a port closest to you)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DOck RElative 1 -1\t\t(Dock to the port relative to your ship by 1 -1)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DOck ABsolute 37 53\t\t(Dock to the port at the location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DOck CLosest PLanet\t\t(Dock at the planet closest to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "DOck CLosest BAse\t\t(Dock at the closes base to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("energy", sndto[1], lenOfString1):
				nme := "Energy"
				syntx := "Energy [<username> | <RElative | ABsolute> <Locy Locx> amount]"
				fctn := "Show your energy or transfer energy between two ships (shields must be down for both ships)."
				optns := "username | RElative | ABsolute"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "ENergy\t\t\t\t(Show your current energy left)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ENergy RElative 1 -1 250\t(Transfer 250 to the ship 1 -1 relative to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ENergy ABsolute 37 53 250\t(Transfer 250 to the ship at location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ENergy Dog -250\t\t\t(Transfer 250 from dog to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ENergy Dog 350\t\t\t(Transfer 350 from your ship to dog)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("gripe", sndto[1], lenOfString1):
				nme := "Gripe"
				syntx := "Gripe [{word ... word}]"
				fctn := "Read or record bugs, comments, suggestion,  etc. in the file gripe.txt, which is periodically reviewed by the implementor."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Gripe"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Gripe Please improve the help system!"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("help", sndto[1], lenOfString1):
				nme := "Help"
				syntx := "Help [command]"
				fctn := "List or describe the legal commands."
				optns := "? or * Shows list of commands"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Help"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Help SUmmary"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("impulse", sndto[1], lenOfString1):
				nme := "Impulse"
				syntx := "Impulse [Absolute|Relative|Computed] <Locy Locx> <Coalition | EMpire | NEutral | ARcheron | FRiendly | ENemy | Ship | Planet | Star | Base | Black Hole>"
				fctn := "Move using impulse engines."
				optns := "Absolute | Relative | Computed"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Impulse RElative 1 -1\t\t(Impulse to the location relative to your ship 1 -1)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Impulse ABsolute 37 53\t\t(Impulse to the absolute location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Impulse COmputed Planet\t\t(Impulse to the closest planet to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Impulse COmputed Base\t\t(Impulse to the closest base to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Impulse COmputed Neutral\t(Impulse to the closest neutral object)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("list", sndto[1], lenOfString1):
				nme := "List"
				syntx := "List [ALl | CLosest | CAptured | FRiendly] {options}"
				//				syntx := "List [SHips | BAses | PLanets | POrts | COalition | EMpire | NEutral | ARcheron | FRiendly | ENemy | TArgets | CAptured | ALl | SUmmary | CLosest]"
				fctn := "List various information  about:  ships,  bases, planets, ports, Coalition, Empire, Neutral, Archeron, Friendly, Enemy, TArgets, Captured, All, Summary, Closest."
				//				optns := "CAptured | ALl | SUmmary | CLosest"
				optns := "SUmmary | SHips | BAses | PLanets | STars | POrts | COalition | EMpire | NEutral | ARcheron | TArgets"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "List\t\t(lists everything seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List ALl\t\t\t(lists everything seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List Closest\t\t\t(lists closest thing seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List CAptured\t\t\t(lists captured planets for your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List ENemy\t\t\t(shows all enemy planets, bases, ships seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List FRiendly\t\t\t(shows objects on your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List SUmmary\t\t\t(lists everything seen by your side with a summary at end)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List SHips\t\t\t(lists ships seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List BAses\t\t\t(lists bases seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List PLanets\t\t\t(lists planets seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List STars\t\t\t(lists stars seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List POrts\t\t\t(lists ports seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List COalition\t\t\t(lists coalition seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List EMpire\t\t\t(lists empire seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List NEutral\t\t\t(lists neutral seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List ARcheron\t\t\t(lists archeron seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List TArgets\t\t\t(lists targets seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List ALl {ships,bases,planets, ports, coalition, empire, neutral, archeron, targets, stars, friendly, enemy, captured}"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "\t\t\t\t(lists objects seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List CLosest {all, ships, stars, bases, planets, ports, coalition, empire, neutral, archeron, enemy, targets, captured}"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "\t\t\t\t(lists closest objects seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List CAptured {all, planets}\t(lists closest objects seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "List FRiendly {all, ships, bases, planets, ports, targets, stars}"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "\t\t\t\t(lists closest objects seen by your side)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("move", sndto[1], lenOfString1):
				nme := "Move"
				syntx := "Move [Absolute|Relative|Computed] <Locy Locx> <Coalition | EMpire | NEutral | ARcheron | FRiendly | ENemy | Ship | Planet | Star | Base | Black Hole | Object>"
				fctn := "Move using warp engines."
				optns := "Absolute | Relative | Computed"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Move RElative 1 -1\t\t(Move to the location relative to your ship 1 -1)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Move ABsolute 37 53\t\t(Move to the absolute location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Move Computed PLanet\t\t(Use the computer to move to the closest planet)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Move Computed Friend\t\t(Use the computer to move to the closest friendly object)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Move Computed p12\t\t(Use the computer to move to planet #12)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("news", sndto[1], lenOfString1):
				nme := "News"
				syntx := "News"
				fctn := "Tell about any  new  features  or enhancements."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "News"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch(NamePhas, sndto[1], lenOfString2):
				nme := NamePhas
				syntx := "PHasers [Computed | Relative | Absolute] <Locy Locx> [Empire | Neutral | Archeron | Friendly | ENemy | Ship | Planet | Star | Base | Black Hole> <amount>"
				fctn := "Fire phasers at a target.  Amount may be 50-500."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "PHasers 100 RElative 1 -1\t(Fire 100 units at relative location 1 -1)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PHasers ABsolute 37 53\t\t(Fire 200 units at absolute location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PHasers Computed bh26\t\t(Fire 200 units at bh26)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PHasers Computed Neutral\t(Fire 200 units at the closest neutral object)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PHasers Computed Planet\t\t(Fire 200 units at the closest planet)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("planets", sndto[1], lenOfString2):
				nme := "PLanets"
				syntx := "PLanets [ENemy | SUmmary | ALl | CLosest | RElative | ABsolute] location]"
				fctn := "List information on friendly and known enemy and neutral planets."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "PLanets\t\t\t\t(List all known planets with summary)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PLanets ALl\t\t\t(List all known planets)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PLanets ENemy\t\t\t(List all known enemy planets)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PLanets SUmmary\t\t\t(List a summary of all known planets)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PLanets CLosest\t\t\t(List the closest known planet)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PLanets RElative 1 -1\t\t(List the known planet relative to your ship by 1 -1)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "PLanets ABsolute 37 53\t\t(List the known planet at the lcoation 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("points", sndto[1], lenOfString2):
				nme := "POints"
				syntx := "POints [{friendly|enemy}]"
				fctn := "List scores breakdown so far."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "POints\t\t\t(Show all points)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "POints Friendly\t\t(Show your side's points"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "POints Enemy\t\t(Show the enemy side's points"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("quit", sndto[1], lenOfString1):
				nme := "Quit"
				syntx := "QUit"
				fctn := "Leave the game and close the connection"
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "QUit"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("radio", sndto[1], lenOfString2):
				nme := "RAdio"
				syntx := "RAdio [<on|off>]"
				fctn := "Turn ship's sub-space radio on or off;  ignore or restore communications from individual ships."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				break

			case testCommandMatch("repair", sndto[1], lenOfString2):
				nme := "REpair"
				syntx := "REpair [radio|shields|tractor|Warp|Impulse|Photon|Phaser|Computer|Life]"
				fctn := "Repair your damaged devices a little."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "REpair\t\t\t(Repair all devices)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "REpair Radio\t\t(Repair the radio)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "REpair Engines\t\t(Repair the engines)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("scan", sndto[1], lenOfString2):
				nme := "SCan"
				syntx := "SCan [Up|Down|Right|Left] [<range>|<vr><hr>] [Warning]"
				fctn := "Display a selected portion of the nearby universe."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "SCan\t\t\t(Default scan)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SCan Up\t\t\t(Scan the quadrant above your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SCan 3 5\t\t(Scan a range 3 vertical and 5 horizontal of your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SCan 3 3 Warning\t(Scan a range 3 vert and horizontal with dangerous space shown as !)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("set", sndto[1], lenOfString2):
				nme := "SEt"
				syntx := "SEt <SIde | PAssword | PRompt | Output | IOfmt | INitial | SCripts> <PASSWORD | PROMPT | Medium | Long | Short | ABsolute | RElative | BOth | Empire | Coalition | Command line>"
				fctn := "Set various input, output, password, prompting and initial command defaults."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "SEt\t\t\t\t(Show current settings)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt SIde Empire\t\t\t(Set your side to empire (available using set initial only)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt PRompt Bigdaddy\t\t(Set your prompt to Bigdaddy>)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt Output Long\t\t\t(Set your output to the long setting"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt Output Short\t\t(Set your output to the short setting)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt IOfmt RElative\t\t(Set your input and output format to relative only)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt IOfmt ABsolute\t\t(Set your input and output format to absolute only)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt IOfmt Both\t\t\t(Set your I/O format to both absolute and relative, input is relative)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt INitial set side emp\t(Set the initial command to a set your side to empire)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt SCripts\t\t\t(this shows what scripts are defined currently)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SEt 0 mo 5 5\\mo 5 5\t\t(create a script, run by typing 0)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("shields", sndto[1], lenOfString2):
				nme := "SHields"
				syntx := "SHields [{Up|Down|<Energy Transfer number>]"
				fctn := "Transfer energy to  or  from  your  shields; raise or lower your shields."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "SHields Up\t\t\t(Raise your shields)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SHields Down\t\t\t(Lower your shields)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SHields Transfer 100\t\t(Transfer 100 units of energy from the shields to the ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SHields Transfer -150\t\t(Transfer 150 units of energy from the ship to your shields)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("srscan", sndto[1], lenOfString2):
				nme := "SRscan"
				syntx := "SRscan"
				fctn := "Display the galaxy with a default range of  7 sectors (1 greater than the maximum warp factor)."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "srscan"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("status", sndto[1], lenOfString2):
				nme := "STatus"
				syntx := "STatus [{RAdio|SHields|ENgines|STardate|TOrp}] | username | RElative | ABsolute] location]"
				fctn := "List a objects current  status  and  supply levels."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "STatus\t\t\t\t(Show the status of your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "STatus P54\t\t\t(Show the status of planet #54 if it is within scanner range)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "STatus SHields\t\t\t(Show the status of your shields)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "STatus RElative 1 -1\t\t(how the status of the object relative 1 -1 from your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "STatus ABsolute 37 53\t\t(Show the status of the object at location 37 53 if within scanner range)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch(NameSummary, sndto[1], lenOfString2):
				nme := "SUmmary"
				syntx := "SUmmary[{SHips, Bases, Planets, Stars, Black Holes}]"
				fctn := "List various information  on  ships,  bases, and planets."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "SUmmary\t\t\t(Show a summary of all objects in the universe)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary SHips\t\t(Show a summary of all ships in the universe)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary Bases\t\t(Show a summary of all bases in the universe)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary Planets\t\t\t(Show a summary of all planets in the universe)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary Black\t\t\t(Show a summary of all back holes in the universe)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("targets", sndto[1], lenOfString2):
				nme := "TArgets"
				syntx := "TArgets"
				fctn := "List  targets  (enemies  within  range)  and their current locations."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TArgets"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("tell", sndto[1], lenOfString2):
				nme := "TEll"
				syntx := "TEll <COalition | EMpire | NEutral | ARcheron | ALl | FRiendlies | ENemies | Username> {word ... word}"
				fctn := "Send  messages  to  other  ships  using   the sub-space radio."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TEll hsn lets attack 57-33"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TEll COalition I saw a archeron at 33-17"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TEll Empire You all suck!"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TEll Neutral Prepare to die!"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("time", sndto[1], lenOfString2):
				nme := "TIme"
				syntx := "TIme"
				fctn := "List information on run time and elapsed time."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TIme"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("torpedoes", sndto[1], lenOfString2):
				nme := NameTorp
				syntx := "TOrpedoes #burst [Absolute|Relative|Computed] [<Locy Locx> | name]"
				fctn := "Fire photon torpedoes at a target."
				optns := "#burst"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TOrpedoes 1 RElative 2 -3\t(Send 1 torpedo toward the object 2 -3 relative to your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TOrpedoes 2 ABsolute 37 53\t(Send 2 torpedos toward object at location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TOrpedoes 1 Computed p26\t(Send 1 torpedo at the planet #26 (if within range)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("tractor", sndto[1], lenOfString2):
				nme := "TRactor"
				syntx := "TRactor [<on, off> <relative, absolute, object> <Locx Loxy>]"
				fctn := "Use tractor beam to tow friendly ships. Both ships shields must be down."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TRactor on RElative 1 -1\t\t(Tractor beam object relative 1 -1 to you)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TRactor on ABsolute 37 53\t\t(Tractor beam object at the location 37 53)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TRactor off\t\t\t\t(Turn off the tractor beam)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TRactor on hsn\t\t\t\t(Turn on the tractor beam to object hsn)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("type", sndto[1], lenOfString2):
				nme := "TYpe <System | InputOutput | Game>"
				syntx := "TYpe user/side message"
				fctn := "List  current Input/Output, system and game characteristics."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TYpe System\t\t\t(Display system information)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TYpe InputOutput\t\tDisplay the input output settings for your ship)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TYpe Game\t\t\t(Display game information"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("users", sndto[1], lenOfString1):
				nme := "USers"
				syntx := "USers"
				fctn := "List all players currently logged on"
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "USers"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("ADministrator", sndto[1], lenOfString1):
				nme := "ADministrator"
				syntx := "ADministrator <bounce | disable | enable | report | remove | su-on | su-off | bh-off | bh-on | ar-on | ar-off | delete | move | reset | wed | ied | ptd | pbd | sd | cd | lsd | rd | tbd | tsd | start | endgame> <username/object> <Locy Locx><amount><phaser><torp>"
				fctn := "Administrative actions for the system"
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "ADministrator bounce hsn\t\t(Force a user off the system)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator disable hsn\t\t(Disable a user account login)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator report\t\t\t(Report a user or all users accounts)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator su-on\t\t\t(Turn on administrative rights for a user)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator bh-off\t\t\t(Turn off administrative rights for a user)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator delete hsn\t\t(Delete an object from the game)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator move hsn 37 53\t\t(Move an object to a new location in the game)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator reset hsn\t\t\t(Send a password recovery email to a user)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator ar-off\t\t\t(Turn off Archerons in the game)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator ar-on\t\t\t(Turn on Archerons in the game)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator wed hsn 500\t\t(sets warp engine damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator ied hsn 500\t\t(sets impulse engine damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator ptd hsn 500\t\t(sets photon torpedos damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator pbd hsn 500\t\t(sets phaser bank damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator sd hsn 500\t\t(sets shield damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator cd hsn 500\t\t(sets computer damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator lsd hsn 500\t\t(sets life support damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator rd hsn 500\t\t(sets radio damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator tbd hsn 500\t\t(sets tractor beam damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator tsd hsn 500\t\t(sets total ship damage to 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator remove hsn 500\t\t(remove the user from the users db)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator dohit hsn wolf torp 500\t(do a hit from hsn to wolf type: torp amount: 500)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator STart Vmax Hmax Starsmax BHsmax Archeronmax Planetsmax Port#"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator STart 75 75 120 120 5 60 1702-1799\t(all fields required)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "\t\t\t\t\t(Defaults shown except for the port ranges allowed)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "ADministrator ENDGame\t\t\t(Ends the game on the port you are connected to)"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			default:
				sendln(conn, invparm)
				break
			}
		}
	} else {
		invalidcommand = dodefHelp(invalidcommand, conn)
		return invalidcommand
	}
	return invalidcommand
}

//
// Pregame Help command
//
func processPregameHelp(comnd string, invalidcommand bool, conn *Conn) bool {
	sndto := strings.Split(strings.TrimSpace(strings.Trim(comnd, ctlonly)), " ")
	// just asked for help
	if len(sndto) == 2 {
		if strings.ToLower(sndto[1]) == "?" || strings.ToLower(sndto[1]) == "*" || len(sndto) != 2 {
			invalidcommand = dodefPregameHelp(invalidcommand, conn)
		} else {
			switch true {
			case testCommandMatch("activate", sndto[1], lenOfString1):
				nme := "Activate"
				syntx := "Activate"
				fctn := "Enter the game from the pregame."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Activate"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("disconnect", sndto[1], lenOfString1):
				nme := "Disconnect"
				syntx := "Disconnect"
				fctn := "Disconnect from the server"
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Disconnect"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("gripe", sndto[1], lenOfString1):
				nme := "Gripe"
				syntx := "Gripe [{word ... word}]"
				fctn := "Read or record bugs, comments,  suggestion,  etc.   in the file gripe.txt, which is periodically reviewed by the implementor."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Gripe Please improve the help system!"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("help", sndto[1], lenOfString1):
				nme := "Help"
				syntx := "Help [command]"
				fctn := "List or describe the legal commands."
				optns := "? or * Shows list of commands"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Help"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "Help SUmmary"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("login", sndto[1], lenOfString2):
				nme := "LOgin"
				syntx := "LOgin <username> <password>"
				fctn := "Authenticate to the system"
				optns := "? or * Shows list of commands"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "LOgin myusername mypassword"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("news", sndto[1], lenOfString1):
				nme := "News"
				syntx := "News"
				fctn := "Info about any  new  features  or enhancements."
				optns := "? or * Shows list of commands"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "News"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("points", sndto[1], lenOfString1):
				nme := "Points"
				syntx := "Points [{friendly|enemy}]"
				fctn := "List scores breakdown so far."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "POints Friendly"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "POints Enemy"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("quit", sndto[1], lenOfString1):
				nme := "Quit"
				syntx := "Quit."
				fctn := "Leave the game and close the connection"
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Quit"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("recover", sndto[1], lenOfString1):
				nme := "RECover"
				syntx := "RECover <email address>.  You will be sent a recovery password"
				fctn := "Recover your lost password to use the system."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "RECover myemailaddress"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("register", sndto[1], lenOfString3):
				nme := "REGister"
				syntx := "REGister <username> <email address>.  Username <=6 characters, no special characters allowed. You will be sent your password"
				fctn := "Register to use the system."
				optns := "? or * Shows list of commands"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "REGister requestedusername myemailaddress"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch(NameSummary, sndto[1], lenOfString2):
				nme := "SUmmary[{SHips, Bases, Planets, Stars, Black Holes}]"
				syntx := "SUmmary"
				fctn := "List various information  on  ships,  bases, and planets."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "SUmmary"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary SHips"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary Bases"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary Planets"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "SUmmary Black"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("targets", sndto[1], lenOfString2):
				nme := "TArgets"
				syntx := "TArgets [{ships|planets|stars|bases}]"
				fctn := "List  targets  (enemies  within  range)  and their current locations."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TArgets"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TArgets Ships"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TArgets Planets"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TArgets Stars"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				example = "TArgets Bases"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("time", sndto[1], lenOfString2):
				nme := "TIme"
				syntx := "TIme"
				fctn := "List information on run time and elapsed time."
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "TIme"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("users", sndto[1], lenOfString1):
				nme := "Users"
				syntx := "Users"
				fctn := "List all players currently logged on"
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "Users"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			case testCommandMatch("logoff", sndto[1], lenOfString2):
				nme := "LOGOff"
				syntx := "LOGOff"
				fctn := "Log off the system"
				optns := "none"
				invalidcommand = doHelp(invalidcommand, conn, nme, syntx, fctn, optns)
				example := "LOGOff"
				invalidcommand = doExample(invalidcommand, conn, example, syntx, fctn, optns)
				break

			default:
				sendln(conn, invparm)
				break
			}
		}
	} else {
		invalidcommand = dodefPregameHelp(invalidcommand, conn)
		return invalidcommand
	}
	return invalidcommand
}

//
// Impulse command  only allow moves of 1
//
func processImpulse(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB) bool {
	var Nme string
	var Nme1 string
	var Locx int
	var Locx1 int
	var Locy int
	var Locy1 int
	var SeenbyEnemy int
	var Objtype int
	var ShipEnergy int
	var ShldUp int
	var Side int
	var MySide int
	var IOFmt int
	var ImpEngDam int
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var actualMoves int
	var energyE int
	var rslt bool
	var newlocx int
	var newlocy int
	var myTractorOn int
	var myTractorWho string
	var movfromy int
	var movfromx int
	var relx int
	var rely int

	//   Move [Absolute|Relative|Computed] <vpos> <hpos>
	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	if len(prm1) > 4 {
		sendln(conn, "Too many parameters")
		return true
	}
	if len(prm1) < 3 {
		sendln(conn, "Too few parameters")
		return true
	}
	//
	// Get the user's data
	//
	err := objectsdb.QueryRow("select Nme, Locx, Locy, IOFmt, ImpEngDam, Side from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &IOFmt, &ImpEngDam, &MySide)
	//if len(prm1) = 3 then both parms must be int
	if len(prm1) == 3 { //2 parms
		parameter1 = strings.Trim(prm1[1], ctlsp)
		parameter2 = strings.Trim(prm1[2], ctlsp)

	} else { // 3 parms
		parameter1 = strings.Trim(prm1[1], ctlsp)
		parameter2 = strings.Trim(prm1[2], ctlsp)
		parameter3 = strings.Trim(prm1[3], ctlsp)
		if _, err := strconv.Atoi(parameter1); err == nil {
			sendln(conn, invparm)
		} else {
			if _, err := strconv.Atoi(parameter2); err != nil {
				sendln(conn, invparm)
				return true
			}
			if _, err := strconv.Atoi(parameter3); err != nil {
				sendln(conn, invparm)
				return true
			}
		}
	}
	//
	// Was the request for absolute - if so, convert parms to relative
	//
	//		pmx := parameter2
	//		pmy := parameter3
	parm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")
	if testCommandMatch("absolute", parm[1], lenOfString1) == true {
		// force users format to absolute
		IOFmt = IOFmtAbs
		parameter1 = parameter2
		parameter2 = parameter3
	} else {
		if testCommandMatch(NameRelative, parm[1], lenOfString1) == true {
			// force users format to relative
			IOFmt = IOFmtRel
			parameter1 = parameter2
			parameter2 = parameter3
		} else { //computed moves to closest matching object - compute the loc and force it to absolute!
			if testCommandMatch(NameComputed, parm[1], lenOfString1) == true {
				// need to use computer here
				newlocx, newlocy, rslt = docomputed(username, parameter2, objectsdb)
				if rslt {
					parameter1 = strconv.Itoa(newlocx)
					parameter2 = strconv.Itoa(newlocy)
					IOFmt = IOFmtAbs
				}
			}
		}
	}

	pmx, err := strconv.Atoi(parameter1)
	pmy, err := strconv.Atoi(parameter2)

	relx = 0
	rely = 0

	// Ensure move is in range allowed
	if IOFmt == IOFmtAbs {
		relx = pmx - Locx
		rely = pmy - Locy
	} else {
		relx = pmx
		rely = pmy
	}

	// Are they moving out of bounds?
	if Locx+relx > Vmax {
		relx = Vmax - Locx
	}
	if Locy+rely > Hmax {
		rely = Hmax - Locy
	}

	if Locx+relx < 0 {
		relx = -Locx
	}
	if Locy+rely < 0 {
		rely = -Locy
	}

	//
	// Perform the move, checking to see if there is a blockage and checking for drive by hits
	//
	//
	nummoves, incx, incy := doCalcIncMove(relx, rely)

	// Impulse moves can only be 1 and no more!
	if nummoves > 1 {
		nummoves = 1
	}

	if nummoves <= 0 {
		sendln(conn, invparm)
		return false
	}

	// Are engines damaged? Reduce to warp 3 if necessary
	if ImpEngDam > KCRIT {
		sendln(conn, "Impulse engines damaged.")
		return false
	}

	//
	//	Are we being tractored - if so break the tractor beam before the move
	//
	objectsdb.QueryRow("select TractorOn from objects WHERE Nme = ?", username).Scan(&myTractorOn)
	if myTractorOn == 1 {
		endTractorOn(username, conn, objectsdb)
	}

	// loop thru moves checking for blocks,
	actualMoves = 0
	movfromx = Locx
	movfromy = Locy

	for i := 1; i <= nummoves; i++ {
		newlocx := Locx + int(incx*float32(i))
		newlocy := Locy + int(incy*float32(i))
		err = objectsdb.QueryRow("select Nme, Locx, Locy, Objtype, SeenbyEnemy, Side from objects WHERE Locx = ? and Locy = ?", newlocx, newlocy).Scan(&Nme1, &Locx1, &Locy1, &Objtype, &SeenbyEnemy, &Side)
		if err != nil {
			//
			// nobody has seen me - i've moved, only my side should be able to see me
			//
			SeenbyEnemy = MySide
			//
			_, err = objectsdb.Exec("UPDATE objects set Locx = ?, Locy = ?, SeenbyEnemy = ? WHERE Nme = ?", newlocx, newlocy, SeenbyEnemy, username)

			// Increment actualMoves for energy consumption
			actualMoves = actualMoves + 1

			// if your towing someone, recurse into this with it's stuff information, move them to old addr
			objectsdb.QueryRow("select TractorWho from objects WHERE Nme = ?", username).Scan(&myTractorWho)
			// xxxxxxx
			if myTractorWho != "" {
				doObjDrag(myTractorWho, movfromx, movfromy, objectsdb)
				movfromx = newlocx
				movfromy = newlocy
			}

		} else { // If it's a black hole you die, otherwise, no damage (right now)
			if Objtype == TypeBH {
				//				notify(username, conn, newlocx, newlocy, msgObjSwallowedBH, Nme1, newlocx, newlocy, objectsdb, 0, 0)
				//
				// handle killing off tractor beams
				//
				//				endTractor(username, conn, objectsdb)
				//
				//				dbDelobjects(conn, objectsdb, username, pointsdb, usersdb)
				//
				//log quit
				//
				//				log.Print(username + " fell into black hole.")
				//
				//				objUpdateActv(Off, username, objectsdb)
				//				return false
				sendln(conn, "Navigation Officer:  \"Captain, computer has overrided our attempt to enter the event horizon!\"")
				break
			} else {
				sendln(conn, "Navigation Officer:  \"Collision averted, Captain!\"")
				break
			}
		}
	}

	// Get my energy
	err = objectsdb.QueryRow("select ShipEnergy, ShldUp from objects WHERE Nme = ?", username).Scan(&ShipEnergy, &ShldUp)
	//

	//
	// Move was successful, must do the following: status goes green, and charge energy for it.
	// from compuserve:
	// shields down, = 6 x dist x dist = E
	// shields up = E x 2
	// Tractoring = E x 3
	//  **** TO DO: Tractoring takes 3 times energy!
	//
	energyE = actualMoves * actualMoves * 6

	if ShldUp == On {
		objUpdateShipEnergy(ShipEnergy-(energyE*2), username, objectsdb)
	} else {
		objUpdateShipEnergy(ShipEnergy-energyE, username, objectsdb)
	}
	objUpdateStat(StatG, username, objectsdb)

	//  Risk potential impulse engine damage (.125%)
	if rand.Float32() < .125 {
		damage := rand.Intn(int((actualMoves * ImpEngineDamageFactor) + 1))
		sendln(conn, "EEEEERRRRRROOOOOOOMMMMMmmmmm!!")
		send(conn, "Captain, the engines suffered ")
		send(conn, strconv.Itoa(damage))
		sendln(conn, " units of damage.")
		objUpdateImpDam(username, damage, objectsdb)
	}

	//
	doGameDelay(username, nummoves, objectsdb)
	return false
}

//
// Do the default Impulse command
//
func dodefImpulse(invalidcommand bool, conn *Conn) bool {
	invalidcommand = processHelp("help impul", invalidcommand, conn)
	return false
}

//
// Output the list
//
func outputList(Nme string, myIOFmt int, Objtype int, Locx int, myLocx int, Locy int, myLocy int, Shld int, OutputLen int, conn *Conn, objectsdb *sql.DB) {
	//
	// use name addr instead?
	//
	msg := nmeaddr(Nme, myLocx, myLocy, myIOFmt, OutputLen, objectsdb)
	send(conn, msg)
	//
	/*
		send(conn, Nme)
		send(conn, "\t")
		// out absolute,relative or both

		tst := myIOFmt & IOFmtAbs
		if tst != 0 {
			send(conn, "@")
			send(conn, strconv.Itoa(Locx))
			send(conn, "-")
			send(conn, strconv.Itoa(Locy))
			send(conn, "\t")
		}
		// Relative only format?
		tst = myIOFmt & IOFmtRel
		if tst != 0 {
			tmp := Locx - myLocx
			send(conn, strconv.Itoa(tmp))
			send(conn, " ")
			tmp = Locy - myLocy
			send(conn, strconv.Itoa(tmp))
			send(conn, "\t")
		}
	*/
	// if base or ship ut out percent shields
	//	tst1 := Objtype & TypeShip
	//	tst2 := Objtype & TypeBase

	//	if (tst1 != 0) || (tst2 != 0) {
	//		shldpct := Shld / InitShield * 100
	//		send(conn, " ")
	//		send(conn, strconv.Itoa(shldpct))
	//		send(conn, "%%")
	//	}
	// end the line
	sendln(conn, " ")
	return
}

//
// Do the actual query
//
func doListQuery(invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string, qry string, mySide int, myLocx int, myLocy int, myIOFmt int, myOutputLen int) bool {
	// fmt.Println("query:",qry)
	rows, err := objectsdb.Query(qry)
	// fmt.Println("err:",err)
	if err == nil {
		for rows.Next() {
			var Nme string
			var Locx int
			var Locy int
			var Side int
			var SeenbyEnemy int
			var Objtype int
			var Shld int
			rows.Scan(&Nme, &Locx, &Locy, &Side, &SeenbyEnemy, &Objtype, &Shld)
			// Always only output stuff seen by your side!
			tst := SeenbyEnemy & mySide
			if tst != 0 {
				outputList(Nme, myIOFmt, Objtype, Locx, myLocx, Locy, myLocy, Shld, myOutputLen, conn, objectsdb)
			}
		}
		rows.Close()
	}
	sendln(conn, " ")
	return false
}

//
// Find the closest object near me, optional arguements include object type (if 0 ignore)
// Increment counter from 1 to maxscanrng
// Do 4 queries with the ranges = (-counter to +counter) from my x/y coord, and the other coord = (-counter to +counter) from my x/y coord
//
func findClosest(myLocx int, myLocy int, optionalType int, optionalSide int, objectsdb *sql.DB) (int, int, bool) {
	var clLocx int
	var clLocy int
	var clObjtype int
	var clObjside int

	for q := 1; q <= MaxScanRng; q++ {
		// do Top search
		qry := "select Locx, Locy, Objtype, Side from objects where Locx>=" + strconv.Itoa(myLocx-q) + " and Locx<=" + strconv.Itoa(myLocx+q) + " and Locy =" + strconv.Itoa(myLocy-q)
		rows, _ := objectsdb.Query(qry)

		for rows.Next() {
			rows.Scan(&clLocx, &clLocy, &clObjtype, &clObjside)
			// We should always close the rows at the end of the query....if the query finds a matching object to the request
			if (optionalType == 0 || optionalType == clObjtype) && (optionalSide == 0 || optionalSide == clObjside) {
				rows.Close()
				return clLocx, clLocy, true
			}
		}
		// do Right search
		qry = "select Locx, Locy, Objtype, Side from objects where Locy>=" + strconv.Itoa(myLocy-q) + " and Locy<=" + strconv.Itoa(myLocy+q) + " and Locx =" + strconv.Itoa(myLocx-q)
		rows, _ = objectsdb.Query(qry)
		for rows.Next() {
			rows.Scan(&clLocx, &clLocy, &clObjtype, &clObjside)
			// We should always close the rows at the end of the query....if the query finds a matching object to the request
			if (optionalType == 0 || optionalType == clObjtype) && (optionalSide == 0 || optionalSide == clObjside) {
				rows.Close()
				return clLocx, clLocy, true
			}
		}
		// do Bottom search
		qry = "select Locx, Locy, Objtype, Side from objects where Locx>=" + strconv.Itoa(myLocx-q) + " and Locx<=" + strconv.Itoa(myLocx+q) + " and Locy =" + strconv.Itoa(myLocy+q)
		rows, _ = objectsdb.Query(qry)
		for rows.Next() {
			rows.Scan(&clLocx, &clLocy, &clObjtype, &clObjside)
			// We should always close the rows at the end of the query....if the query finds a matching object to the request
			if (optionalType == 0 || optionalType == clObjtype) && (optionalSide == 0 || optionalSide == clObjside) {
				rows.Close()
				return clLocx, clLocy, true
			}
		}
		// do Left search
		qry = "select Locx, Locy, Objtype, Side from objects where Locy>=" + strconv.Itoa(myLocy-q) + " and Locy<=" + strconv.Itoa(myLocy+q) + " and Locx =" + strconv.Itoa(myLocx+q)
		rows, _ = objectsdb.Query(qry)
		for rows.Next() {
			rows.Scan(&clLocx, &clLocy, &clObjtype, &clObjside)
			// We should always close the rows at the end of the query....if the query finds a matching object to the request
			if (optionalType == 0 || optionalType == clObjtype) && (optionalSide == 0 || optionalSide == clObjside) {
				rows.Close()
				return clLocx, clLocy, true
			}
		}
	}
	//
	// If we don't get anything, return false
	//
	return clLocx, clLocy, false
}

//
// List command
// List various information  about:  ships,  bases, planets, ports, Coalition, Empire, Neutral, Archeron, Friendly, Enemy, TArgets, Captured, All, Summary, Closest.
// syntax:
// list
// list <verb> <noun>
// Verbs: All, Closest, Captured, Enemy, Friendly, Summary
// Nouns: ships,  bases, planets, ports, stars, Coalition, Empire, Neutral, Archeron, TArgets
//
func processList(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	//
	// What side am I on
	//
	var mySide int
	var myLocx int
	var myLocy int
	var myIOFmt int
	var myOutputLen int
	var clLocx int
	var clLocy int
	var qry string
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var rslt bool
	var summaryFlagType int // what should I summarize?  Zero = all

	send(conn, "\n")
	// get my info
	objectsdb.QueryRow("select Side, Locx, Locy, IOFmt, OutputLen from objects where Nme=?", username).Scan(&mySide, &myLocx, &myLocy, &myIOFmt, &myOutputLen)
	//
	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")

	//
	// Default flag is to show all summaries
	//
	summaryFlagType = 0
	//
	// First parm
	parameter1 = strings.Trim(prm1[1], ctlsp)
	// Second parm
	if len(prm1) == 3 {
		parameter2 = strings.Trim(prm1[2], ctlsp)
	} else {
		parameter2 = ""
	}
	// third parm
	if len(prm1) == 4 {
		parameter2 = strings.Trim(prm1[2], ctlsp)
		parameter3 = strings.Trim(prm1[3], ctlsp)
	} else {
		parameter3 = ""
	}
	// fmt.Println("parameter1=",parameter1," parameter2=",parameter2, " parameter3=",parameter3)
	// if 4 parms, 1st must be ca/all/su/cl and last must be su or error
	if len(prm1) == 4 {
		if (testCommandMatch(NameCaptured, parameter1, lenOfString1) == false) && (testCommandMatch(NameAll, parameter1, lenOfString1) == false) && (testCommandMatch(NameSummary, parameter1, lenOfString1) == false) && (testCommandMatch(NameClosest, parameter1, lenOfString1) == false) && (testCommandMatch(NameSummary, parameter3, lenOfString1) == false) {
			sendln(conn, invparm)
			return true
		}
	}

	switch true {
	case ((len(prm1) == 2) || (len(prm1) == 3 && testCommandMatch(NameSummary, parameter2, lenOfString1))): // 1 parms - can be any keyword

		switch true {
		case testCommandMatch(NameAll, parameter1, lenOfString1):
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype != " + strconv.Itoa(TypeStar) + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameClosest, parameter1, lenOfString1):
			//
			// determine closest object to me
			//
			clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, 0, objectsdb)
			if rslt == true {
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
			}

		case testCommandMatch(NameCaptured, parameter1, lenOfString1):
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + strconv.Itoa(TypePlanet) + " and Side != " + strconv.Itoa(mySide) + " and Side != " + strconv.Itoa(SideNeutral)
			rslt = true

		case testCommandMatch(NameEnemy, parameter1, lenOfString1):
			theside := strconv.Itoa(mySide)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side <> " + theside + " and Objtype != " + strconv.Itoa(TypeStar) + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameFriendly, parameter1, lenOfString1):
			theside := strconv.Itoa(mySide)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + theside + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameSummary, parameter1, lenOfString1):
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype != " + strconv.Itoa(TypeStar) + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameShips, parameter1, lenOfString1):
			obty := strconv.Itoa(TypeShip)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + obty + " ORDER by Nme"
			rslt = true
			summaryFlagType = summaryFlagType | TypeShip

		case testCommandMatch(NameBases, parameter1, lenOfString1):
			obty := strconv.Itoa(TypeBase)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + obty + " ORDER by Nme"
			rslt = true
			summaryFlagType = summaryFlagType | TypeBase

		case testCommandMatch(NamePlanets, parameter1, lenOfString1):
			obty := strconv.Itoa(TypePlanet)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + obty + " ORDER by Nme"
			rslt = true
			summaryFlagType = summaryFlagType | TypePlanet

		case testCommandMatch(NameStars, parameter1, lenOfString1):
			obty := strconv.Itoa(TypeStar)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + obty + " ORDER by Nme"
			rslt = true
			summaryFlagType = summaryFlagType | TypeStar

		case testCommandMatch(NamePorts, parameter1, lenOfString1):
			obty := "(Objtype = " + strconv.Itoa(TypePlanet) + " OR Objtype = " + strconv.Itoa(TypeBase) + ")"
			sde := strconv.Itoa(mySide)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where " + obty + " and side = " + sde + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameCoalition, parameter1, lenOfString1):
			theside := strconv.Itoa(SideCoalition)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + theside + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameEmpire, parameter1, lenOfString1):
			theside := strconv.Itoa(SideEmpire)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + theside + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameNeutral, parameter1, lenOfString1):
			theside := strconv.Itoa(SideNeutral)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + theside + " ORDER by Nme"
			rslt = true

		case testCommandMatch(NameArcheron, parameter1, lenOfString1):
			theside := strconv.Itoa(SideArcheron)
			qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + theside + " ORDER by Nme"
			rslt = true

			//
		case testCommandMatch(NameTargets, parameter1, lenOfString1):
			// simular to scan
			headerlow := myLocy - MaxScanRng
			headerhigh := myLocy + MaxScanRng
			if headerlow < 0 {
				headerlow = 0
			}
			if headerhigh > Vmax {
				headerhigh = Vmax
			}
			rowlow := myLocx - MaxScanRng
			rowhigh := myLocx + MaxScanRng
			if rowlow < 0 {
				rowlow = 0
			}
			if rowhigh > Hmax {
				rowhigh = Hmax
			}
			var Nme string
			var SeenbyEnemy int
			var Side int
			var Objtype int
			var Locx int
			var Locy int
			var Shld int
			for r1 := rowhigh; r1 >= rowlow; r1-- {
				for h1 := headerlow; h1 <= headerhigh; h1++ {
					err := objectsdb.QueryRow("select Nme, Locx, Locy, Objtype, Side, SeenbyEnemy, Shld from objects WHERE Locx = ? and Locy = ? and Side <> ? and (Objtype = ? or Objtype = ? or Objtype = ?);", r1, h1, mySide, TypePlanet, TypeShip, TypeBase).Scan(&Nme, &Locx, &Locy, &Objtype, &Side, &SeenbyEnemy, &Shld)
					if err == nil {
						// Must have seen it
						tst := SeenbyEnemy & mySide
						if tst != 0 {
							outputList(Nme, myIOFmt, Objtype, Locx, myLocx, Locy, myLocy, Shld, shrtTomed(myOutputLen), conn, objectsdb)
						}
					}
				}
			}
			rslt = false
			return true

		default:
			sendln(conn, invparm)
			return true

		}

	case ((len(prm1) == 3) || (len(prm1) == 4 && testCommandMatch(NameSummary, parameter3, lenOfString1))): // 2 parms - can be any keyword

		// verb noun - build a query and run it
		qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects "

		//
		// set rslt to true
		//
		rslt = true
		switch true {
		case testCommandMatch(NameAll, parameter1, lenOfString1):
			switch true {
			case testCommandMatch(NameShips, parameter2, lenOfString1):
				qry = qry + "Where Objtype = " + strconv.Itoa(TypeShip)
				rslt = true
				summaryFlagType = summaryFlagType | TypeShip

			case testCommandMatch(NameBases, parameter2, lenOfString1):
				qry = qry + "Where Objtype = " + strconv.Itoa(TypeBase)
				rslt = true
				summaryFlagType = summaryFlagType | TypeBase

			case testCommandMatch(NamePlanets, parameter2, lenOfString2):
				qry = qry + "Where Objtype = " + strconv.Itoa(TypePlanet)
				rslt = true
				summaryFlagType = summaryFlagType | TypePlanet

			case testCommandMatch(NamePorts, parameter2, lenOfString2):
				obty := "(Objtype = " + strconv.Itoa(TypePlanet) + " OR Objtype = " + strconv.Itoa(TypeBase) + ")"
				sde := strconv.Itoa(mySide)
				qry = qry + "where " + obty + " and side = " + sde
				rslt = true

			case testCommandMatch(NameCoalition, parameter2, lenOfString1):
				qry = qry + "Where Side = " + strconv.Itoa(SideCoalition)
				rslt = true

			case testCommandMatch(NameEmpire, parameter2, lenOfString1):
				qry = qry + "Where Side = " + strconv.Itoa(SideEmpire)
				rslt = true

			case testCommandMatch(NameNeutral, parameter2, lenOfString1):
				qry = qry + "Where Side = " + strconv.Itoa(SideNeutral)
				rslt = true

			case testCommandMatch(NameArcheron, parameter2, lenOfString1):
				qry = qry + "Where Side = " + strconv.Itoa(SideArcheron)
				rslt = true

			case testCommandMatch(NameTargets, parameter2, lenOfString1):
				theside := strconv.Itoa(mySide)
				qry = qry + "Where Side <> " + theside
				rslt = true

			case testCommandMatch(NameStars, parameter2, lenOfString1):
				qry = qry + "Where Objtype = " + strconv.Itoa(TypeStar)
				rslt = true
				summaryFlagType = summaryFlagType | TypeStar

			case testCommandMatch(NameFriendly, parameter2, lenOfString1):
				theside := strconv.Itoa(mySide)
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + theside + " ORDER by Nme"
				rslt = true

			case testCommandMatch(NameEnemy, parameter2, lenOfString1):
				theside := strconv.Itoa(mySide)
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side <> " + theside + " ORDER by Nme"
				rslt = true

			case testCommandMatch(NameCaptured, parameter2, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + strconv.Itoa(TypePlanet) + " and Side != " + strconv.Itoa(mySide) + " and Side != " + strconv.Itoa(SideNeutral)
				rslt = true

			default:
				sendln(conn, invparm)
				return false
			}

		case testCommandMatch(NameClosest, parameter1, lenOfString1):
			switch true {
			case testCommandMatch(NameAll, parameter2, lenOfString1):
				//
				// determine closest object to me
				//
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, 0, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}

			case testCommandMatch(NameShips, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeShip, 0, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
					summaryFlagType = summaryFlagType | TypeShip
				}

			case testCommandMatch(NameStars, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeStar, 0, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
					summaryFlagType = summaryFlagType | TypeStar
				}

			case testCommandMatch(NameBases, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeBase, 0, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
					summaryFlagType = summaryFlagType | TypeBase
				}

			case testCommandMatch(NamePlanets, parameter2, lenOfString2):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypePlanet, 0, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
					summaryFlagType = summaryFlagType | TypePlanet
				}

			case testCommandMatch(NamePorts, parameter2, lenOfString2):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypePlanet, mySide, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				} else {
					clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypeBase, mySide, objectsdb)
					if rslt == true {
						qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
					}
				}
			case testCommandMatch(NameCoalition, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideCoalition, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}
			case testCommandMatch(NameEmpire, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideEmpire, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}

			case testCommandMatch(NameNeutral, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideNeutral, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}

			case testCommandMatch(NameArcheron, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideArcheron, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}
			case testCommandMatch(NameFriendly, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, mySide, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}

			case testCommandMatch(NameEnemy, parameter2, lenOfString1):
				if mySide == SideEmpire {
					clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideCoalition, objectsdb)
				} else {
					clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideEmpire, objectsdb)
				}
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}

			case testCommandMatch(NameTargets, parameter2, lenOfString1):
				if mySide == SideEmpire {
					clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideCoalition, objectsdb)
				} else {
					clLocx, clLocy, rslt = findClosest(myLocx, myLocy, 0, SideEmpire, objectsdb)
				}
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				}

			case testCommandMatch(NameCaptured, parameter2, lenOfString1):
				clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypePlanet, SideCoalition, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
				} else {
					clLocx, clLocy, rslt = findClosest(myLocx, myLocy, TypePlanet, SideCoalition, objectsdb)
					if rslt == true {
						qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy)
					}
				}

			default:
				sendln(conn, invparm)
				return false

			}

		case testCommandMatch(NameCaptured, parameter1, lenOfString1):
			switch true {
			case testCommandMatch(NamePlanets, parameter2, lenOfString1):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + strconv.Itoa(TypePlanet) + " and Side != " + strconv.Itoa(mySide) + " and Side != " + strconv.Itoa(SideNeutral)
				} else {
					sendln(conn, invparm)
					return false
				}

			case testCommandMatch(NameEnemy, parameter1, lenOfString1):
				fmt.Printf("got into captured enemy")

			case testCommandMatch(NameAll, parameter2, lenOfString1):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide)
				}

			case testCommandMatch(NameShips, parameter2, lenOfString1):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeShip)
					summaryFlagType = summaryFlagType | TypeShip
				}

			case testCommandMatch(NameBases, parameter2, lenOfString1):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeBase)
					summaryFlagType = summaryFlagType | TypeBase
				}

			case testCommandMatch(NamePlanets, parameter2, lenOfString2):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypePlanet)
					summaryFlagType = summaryFlagType | TypePlanet
				}

			case testCommandMatch(NamePorts, parameter2, lenOfString2):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide) + " and (Objtype = " + strconv.Itoa(TypePlanet) + " or Objtype =" + strconv.Itoa(TypeBase) + ")"
				rslt = true

			case testCommandMatch(NameTargets, parameter2, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide)
				rslt = true

			case testCommandMatch(NameStars, parameter2, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeStar)
				rslt = true
				summaryFlagType = summaryFlagType | TypeStar

			default:
				sendln(conn, invparm)
				return false
			}

		case testCommandMatch(NameFriendly, parameter1, lenOfString1):
			switch true {
			case testCommandMatch(NameAll, parameter2, lenOfString1):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide)
				}

			case testCommandMatch(NameShips, parameter2, lenOfString2):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeShip)
					summaryFlagType = summaryFlagType | TypeShip
				}

			case testCommandMatch(NameBases, parameter2, lenOfString1):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeBase)
					summaryFlagType = summaryFlagType | TypeBase
				}

			case testCommandMatch(NamePlanets, parameter2, lenOfString2):
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypePlanet)
				}

			case testCommandMatch(NamePorts, parameter2, lenOfString2):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide) + " and (Objtype = " + strconv.Itoa(TypePlanet) + " or Objtype =" + strconv.Itoa(TypeBase) + ")"
				rslt = true

			case testCommandMatch(NameTargets, parameter2, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide)
				rslt = true

			case testCommandMatch(NameStars, parameter2, lenOfString2):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side = " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypeStar)
				rslt = true
				summaryFlagType = summaryFlagType | TypeStar

			default:
				sendln(conn, invparm)
				return false
			}
		}
		//			return false

		//other cases

	default:
		sendln(conn, invparm)
		return false
	}

	if rslt == true {
		doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
		if testCommandMatch(NameSummary, parameter1, lenOfString1) || testCommandMatch(NameSummary, parameter2, lenOfString1) || testCommandMatch(NameSummary, parameter3, lenOfString1) {
			switch true {
			case (summaryFlagType & TypeShip) != 0:
				sendln(conn, " ")
				doSummaryShips(conn, objectsdb)

			case (summaryFlagType & TypeBase) != 0:
				sendln(conn, " ")
				doSummaryBases(conn, objectsdb)

			case (summaryFlagType & TypePlanet) != 0:
				sendln(conn, " ")
				doSummaryPlanets(conn, objectsdb)

			case (summaryFlagType & TypeStar) != 0:
				sendln(conn, " ")
				doSummaryStars(conn, objectsdb)

			case summaryFlagType == 0:
				sendln(conn, " ")
				doSummaryShips(conn, objectsdb)
				sendln(conn, " ")
				doSummaryBases(conn, objectsdb)
				sendln(conn, " ")
				doSummaryPlanets(conn, objectsdb)
				sendln(conn, " ")
				doSummaryStars(conn, objectsdb)
			}
		}
	} else {
		sendln(conn, "No results")
		return false
	}

	return false
}

//
// Do the default List command
//
func dodefList(invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	//
	// What side am I on
	//
	var mySide int
	var myLocx int
	var myLocy int
	var myIOFmt int
	var myOutputLen int
	objectsdb.QueryRow("select Side, Locx, Locy, IOFmt, OutputLen from objects where Nme=?", username).Scan(&mySide, &myLocx, &myLocy, &myIOFmt, &myOutputLen)

	qry := "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype != " + strconv.Itoa(TypeStar) + " ORDER by Nme"

	//	old qry := "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects ORDER by Nme"
	doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
	return false
}

//
// hashing code
//
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

//
// Recover command = Recover email address
//
func processRecover(comnd string, invalidcommand bool, conn *Conn, netaddr string, usersdb *sql.DB) {
	var dateSent time.Time
	var username string
	//
	// Regex for parsing input for bad stuff
	//
	//	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	rem := regexp.MustCompile("^[a-zA-Z0-9_@.]*$")
	//
	prm := strings.Trim(comnd, ctlonly)
	sliceemail := strings.SplitAfter(prm, " ")
	//
	// right number of parms?
	//
	if len(sliceemail) != 2 {
		sendln(conn, invparm)
		return
	}
	//
	email := strings.Trim(sliceemail[1], ctlsp)
	if rem.MatchString(email) == false {
		send(conn, "Email contains invalid character!\n")
		return
	}

	//
	// Valid email?
	//
	err := usersdb.QueryRow("select name, RecoveryDateSent from users where mailAddr=?", email).Scan(&username, &dateSent)
	if err != nil {
		sendln(conn, "Invalid email address, or you haven't registered yet. Use the register command to do that prior to logging in.")
		username = ""
		return
	}

	// generate reandom pw
	pw := doRandomPassword()

	//
	// Put it into db
	//
	tx, err := usersdb.Begin()
	_, err = usersdb.Exec("UPDATE users SET pswd =  ?, RecoveryDateSent = ? where name = ?", GetMD5Hash(pw), time.Now().Format("02-Jan-06"), username)
	tx.Commit()

	//
	// Send email with password to user
	//
	body := "To: " + email + "\r\nSubject: Decwars password recovery\r\n\r\n" +
		"You are being sent this because a password reset was received from the Decwars.com game.\r\nYour username is: " + username + "\r\nYour password is:" + pw + "\r\nYou can change your password after entering the game with the SET PAssword command.\r\nPlay decwars at http://decwars.com or telnet decwars.com 1701"
	auth := smtp.PlainAuth("", "decwars@gmail.com", "yourpassword here", "smtp.gmail.com")
	//	err1 := smtp.SendMail("smtp.gmail.com:587", auth, "decwars@gmail.com", []string{email}, []byte(body))
	smtp.SendMail("smtp.gmail.com:587", auth, "decwars@gmail.com", []string{email}, []byte(body))
	// fmt.Println("Error on sending email:", err1)
	//
	// end emailing registration
	//
	send(conn, "An recovery email with your new password has been sent to: ")
	sendln(conn, email)
	return
}

//
// Register command
//
func processRegister(comnd string, invalidcommand bool, conn *Conn, netaddr string, usersdb *sql.DB) bool {
	var disab int
	// registering users are not disabled
	disab = 0
	//
	// Regex for parsing input for bad stuff
	//
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	// orig	rem := regexp.MustCompile("^[a-zA-Z0-9_@.]*$")
	rem := regexp.MustCompile("^([a-zA-Z0-9_\\-\\.]+)@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.)|(([a-zA-Z0-9\\-]+\\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\\]?)$")
	//
	prm := strings.Trim(comnd, ctlonly)
	sliceusername := strings.SplitAfter(prm, " ")
	//
	// right number of parms?
	//
	if len(sliceusername) != 3 {
		sendln(conn, invparm)
		return true
	}
	//
	// force username to lowercase
	//
	username := strings.ToLower(strings.Trim(sliceusername[1], ctlsp))
	email := strings.Trim(sliceusername[2], ctlsp)
	if re.MatchString(username) == false {
		send(conn, "Username contains invalid character!\n")
		return false
	}
	if rem.MatchString(email) == false {
		send(conn, "Email contains invalid character!\n")
		return false
	}

	if len(username) > Usernamemax {
		send(conn, "Usernames may only be up to 6 characters long.\n")
		return false
	}
	
	fmt.Println("len is: ",len(username))
	
	fmt.Println(" first letter:",username[0])
	
	// temp here - generate reandom pw - eventually will be email processing
	pw := doRandomPassword()

	// ******************* take out after testing *************************
	// pw = "a"
	/*
		 * not needed here
		 //
		// Random location for user (needs to check for stomping)
		//
		v := strconv.Itoa(rand.Intn(Vmax))
		h := strconv.Itoa(rand.Intn(Hmax))
	*/
	//
	// Put it into db
	//
	tx, err := usersdb.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO users(name, mailAddr, pswd, addr, disabled, tme, DamToBH, DamToBases, DamToShip, DamToStars, DamToPlanets, NumOfShips, NumOfStarDates, RecoveryDateSent, SuperUser, CumlativeWins) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, email, GetMD5Hash(pw), conn.RemoteAddr().String(), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	if err != nil {
		sendln(conn, "Username has already been taken, please retry or use the recover command.")
		return false
	}
	tx.Commit()

	//
	// Send email with password to user
	//
	body := "To: " + email + "\r\nSubject: Decwars registration\r\n\r\n" +
		"Welcome to Decwars, your password is:" + pw + "\r\nYou may play decwars at decwars.com"

	auth := smtp.PlainAuth("", "decwars", "your password here", "smtp.gmail.com")
	smtp.SendMail("smtp.gmail.com:587", auth, "decwars@gmail.com", []string{email}, []byte(body))
	//
	// end emailing registration
	//
	send(conn, "An email with your password has been sent to: ")
	sendln(conn, email)
	return false
}

//
// Login command
//
func processLogin(comnd string, invalidcommand bool, conn *Conn, netaddr string, username string, usersdb *sql.DB) (bool, string) {
	//
	// Check to see if this connection is already logged in
	//
	if username != "" {
		sendln(conn, "You must log out prior to logging back in")
		return false, username
	}
	//
	// Regex for parsing input for bad stuff
	//
	//	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	passw := ""
	//
	prm := strings.Trim(comnd, ctlonly)
	sliceusername := strings.SplitAfter(prm, " ")
	//
	// Is the number of parms correct?  Must be 2
	//
	if len(sliceusername) != 3 {
		sendln(conn, invparm)
		return false, username
	}
	username = strings.Trim(sliceusername[1], ctlsp)
	if len(username) > Usernamemax {
		prm = ""
		sendln(conn, "Usernames may only be up to 6 characters long.")
		return false, username
	}
	//
	if len(sliceusername) >= 3 {
		passw = strings.Trim(sliceusername[2], ctlonly)
	} else {
		passw = ""
	}

	// check to see if username exists
	//
	_, err := usersdb.Begin()
	if err != nil {
		log.Fatal(err)
	}
	nme := ""
	pswd := ""
	tme := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	disb := 0

	//	addr := ""
	//	usrconn := conn
	err = usersdb.QueryRow("select name, pswd, tme, disabled from users where name=? and pswd=?", username, GetMD5Hash(passw)).Scan(&nme, &pswd, &tme, &disb)
/*
	if err != nil {
		sendln(conn, "Invalid username/password, or you haven't registered yet. Use the register command to do that prior to logging in.")
		username = ""
		return false, username
	}
*/

	// check to see if the account is disabled
	if disb == 1 {
		sendln(conn, "Account disabled. Contact decwars@gmail.com")
		username = ""
		return false, username
	}
	// check to see if they are already logged in, check the conmap
	_, exists := Conmap[username]
	if exists {
		sendln(conn, "That user is already logged in.")
		username = ""
		return false, username
	}

	// Ok to log them in
	
fmt.Println("*** *** *** username=",username, " conn:",conn, " netaddr:",netaddr)
	
	Conmap[username] = Constr{username, conn, netaddr}

	//  incidate they are logged on
	sendln(conn, "Logged into pregame, to play type ACtivate")
	return true, username
}

//
// Logoff command
//
func processLogoff(username string, invalidcommand bool, conn *Conn, netaddr string) (bool, string) {
	if username != "" {
		dbDelete(username)
		sendln(conn, "Successfully logged off")
	} else {
		sendln(conn, "Not logged on")
		return false, username
	}
	return false, ""
}

//
// Drag an object - skip lots of checks etc - need to fix seenby though
//
func doObjDrag(myTractorWho string, movfromx int, movfromy int, objectsdb *sql.DB) {
	var Locx int
	var Locy int
	var AnotherTractorWho string

	objectsdb.QueryRow("select Locx, Locy, TractorWho from objects WHERE Nme = ?", myTractorWho).Scan(&Locx, &Locy, &AnotherTractorWho)

	objectsdb.Exec("UPDATE objects set Locx = ?, Locy = ? WHERE Nme = ?", movfromx, movfromy, myTractorWho)

	// ok, does this object have it's tractor on? If so, you know what to do...
	if AnotherTractorWho != "" {
		doObjDrag(AnotherTractorWho, Locx, Locy, objectsdb)
	}
}

//
// Move command
//
func processMove(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB, recurse bool, pointsdb *sql.DB, usersdb *sql.DB) bool {
	var Nme string
	var Nme1 string
	var Locx int
	var Locx1 int
	var Locy int
	var Locy1 int
	var SeenbyEnemy int
	var Objtype int
	var ShipEnergy int
	var ShldUp int
	var Side int
	var MySide int
	var IOFmt int
	var WarpEngDam int
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var actualMoves int
	var energyE int
	var rslt bool
	var newlocx int
	var newlocy int
	var myTractorOn int
	var myTractorWho string
	var myDockFlag string
	var portLocx int
	var portLocy int

	//   Move [Absolute|Relative|Computed] <vpos> <hpos>
	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	if len(prm1) > 4 {
		sendln(conn, "Too many parameters")
		return true
	}
	if len(prm1) < 3 {
		sendln(conn, "Too few parameters")
		return true
	}

	//
	// Get the user's data
	//
	err := objectsdb.QueryRow("select Nme, Locx, Locy, IOFmt, WarpEngDam, Side from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &IOFmt, &WarpEngDam, &MySide)
	//
	//if len(prm1) = 3 then both parms must be int
	if len(prm1) == 3 { //2 parms
		parameter1 = strings.Trim(prm1[1], ctlsp)
		parameter2 = strings.Trim(prm1[2], ctlsp)

	} else { // 3 parms
		parameter1 = strings.Trim(prm1[1], ctlsp)
		parameter2 = strings.Trim(prm1[2], ctlsp)
		parameter3 = strings.Trim(prm1[3], ctlsp)
		if _, err := strconv.Atoi(parameter1); err == nil {
			sendln(conn, invparm)
		} else {
			if _, err := strconv.Atoi(parameter2); err != nil {
				sendln(conn, invparm)
				return true
			}
			if _, err := strconv.Atoi(parameter3); err != nil {
				sendln(conn, invparm)
				return true
			}
		}
	}

	//
	// Was the request for absolute - if so, convert parms to relative
	//
	//		pmx := parameter2
	//		pmy := parameter3
	parm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")
	if testCommandMatch("absolute", parm[1], lenOfString1) == true {
		// force users format to absolute
		IOFmt = IOFmtAbs
		parameter1 = parameter2
		parameter2 = parameter3
	} else {
		if testCommandMatch(NameRelative, parm[1], lenOfString1) == true {
			// force users format to relative
			IOFmt = IOFmtRel
			parameter1 = parameter2
			parameter2 = parameter3
		} else { //computed moves to closest matching object - compute the loc and force it to absolute!
			if testCommandMatch(NameComputed, parm[1], lenOfString1) == true {
				// need to use computer here
				newlocx, newlocy, rslt = docomputed(username, parameter2, objectsdb)
				if rslt {
					parameter1 = strconv.Itoa(newlocx)
					parameter2 = strconv.Itoa(newlocy)
					IOFmt = IOFmtAbs
				}
			}
		}
	}

	pmx, err := strconv.Atoi(parameter1)
	pmy, err := strconv.Atoi(parameter2)

	relx := 0
	rely := 0

	// Ensure move is in range allowed
	if IOFmt == IOFmtAbs {
		relx = pmx - Locx
		rely = pmy - Locy
	} else {
		relx = pmx
		rely = pmy
	}

	// Are they moving out of bounds?
	if Locx+relx > Vmax {
		relx = Vmax - Locx
	}
	if Locy+rely > Hmax {
		rely = Hmax - Locy
	}

	if Locx+relx < 0 {
		relx = -Locx
	}
	if Locy+rely < 0 {
		rely = -Locy
	}

	//
	// Perform the move, checking to see if there is a blockage and checking for drive by hits
	// if computed, ship will crash so the number can be 1 greater than movebig!
	//
	nummoves, incx, incy := doCalcIncMove(relx, rely)
	if testCommandMatch(NameComputed, parm[1], lenOfString1) == true {
		if nummoves > MaxWarpFactor+1 {
			sendln(conn, movetbig)
			return false
		}
	} else {
		if nummoves > MaxWarpFactor {
			sendln(conn, movetbig)
			return false
		}
	}

	if nummoves <= 0 {
		sendln(conn, invparm)
		return false
	}

	if recurse == false {
		// Are engines damaged? Reduce to warp 3 if necessary
		if WarpEngDam > KCRIT {
			sendln(conn, "Warp engines damaged.")
			return false
		}

		if WarpEngDam > 0 { //this was maxdam
			if nummoves > 3 {
				nummoves = 3
				sendln(conn, "Engines damaged, warp 3 max, adjusting Captain!")
			}
		}

		//
		//	Are we being tractored  - if so break the relationships before the move
		//
		objectsdb.QueryRow("select TractorOn from objects WHERE Nme = ?", username).Scan(&myTractorOn)
		if myTractorOn == 1 {
			endTractorOn(username, conn, objectsdb)
		}
	}

	//
	// loop thru moves checking for blocks
	// need to add tractor support - kill tractor if moved or follow
	actualMoves = 0
	movfromx := Locx
	movfromy := Locy
	for i := 1; i <= nummoves; i++ {
		newlocx := Locx + int(incx*float32(i))
		newlocy := Locy + int(incy*float32(i))
		err = objectsdb.QueryRow("select Nme, Locx, Locy, Objtype, SeenbyEnemy, Side from objects WHERE Locx = ? and Locy = ?", newlocx, newlocy).Scan(&Nme1, &Locx1, &Locy1, &Objtype, &SeenbyEnemy, &Side)
		if err != nil {
			//
			// nobody has seen me - i've moved, only my side should be able to see me
			//
			SeenbyEnemy = MySide
			//
			objectsdb.Exec("UPDATE objects set Locx = ?, Locy = ?, SeenbyEnemy = ? WHERE Nme = ?", newlocx, newlocy, SeenbyEnemy, username)

			// Increment actualMoves for energy consumption
			actualMoves = actualMoves + 1

			// if your towing someone, recurse into this with it's stuff information, move them to old addr
			objectsdb.QueryRow("select TractorWho from objects WHERE Nme = ?", username).Scan(&myTractorWho)
			// xxxxxxx
			if myTractorWho != "" {
				doObjDrag(myTractorWho, movfromx, movfromy, objectsdb)
				movfromx = newlocx
				movfromy = newlocy
			}

		} else { // If it's a black hole you die, otherwise, no damage (right now)
			if Objtype == TypeBH {
				//				notify(username, conn, newlocx, newlocy, msgObjSwallowedBH, Nme1, newlocx, newlocy, objectsdb, 0, 0)

				//
				// handle killing off tractor beams
				//
				//				endTractor(username, conn, objectsdb)

				//				dbDelobjects(conn, objectsdb, username, pointsdb, usersdb)
				//
				//log quit
				//
				//				log.Print(username + " fell into black hole.")
				//
				//				objUpdateActv(Off, username, objectsdb)
				//				return false
				//
				// code inserted to fix moving into black holes like in the real decwar
				//
				sendln(conn, "Navigation Officer:  \"Captain, computer has overrided our attempt to enter the event horizon!\"")
				break
			} else {
				sendln(conn, "Navigation Officer:  \"Collision averted, Captain!\"")
				break
			}
		}
	}

	// Get my energy
	err = objectsdb.QueryRow("select ShipEnergy, ShldUp, Locx, Locy, DockFlag, TractorWho from objects WHERE Nme = ?", username).Scan(&ShipEnergy, &ShldUp, &Locx, &Locy, &myDockFlag, &myTractorWho)
	//
	// If docked, are we adjacent to the port?
	//
	if myDockFlag != "" {
		objectsdb.QueryRow("select Locx, Locy from objects WHERE Nme = ?", myDockFlag).Scan(&portLocx, &portLocy)
		if Abs(Locx-portLocx) > 1 || Abs(Locy-portLocy) > 1 {
			endDock(username, conn, objectsdb)
		}
	}

	//
	// Move was successful, must do the following: status goes green, and charge energy for it.
	// from compuserve:
	// shields down, = 6 x dist x dist = E
	// shields up = E x 2
	// Tractoring = E x 3
	//  **** TO DO: Tractoring takes 3 times energy!
	//
	energyE = actualMoves * actualMoves * 6

	if ShldUp == On {
		energyE = energyE * 2
	}

	if myTractorWho != "" {
		energyE = energyE * 3
	}

	objUpdateShipEnergy(ShipEnergy-energyE, username, objectsdb)

	objUpdateStat(StatG, username, objectsdb)

	//  Warp factors 5 and 6 risk potential warp engine damage (25%)
	if actualMoves > 4 {
		if rand.Float32() <= .25 {
			damage := rand.Intn(int((actualMoves * WarpEngineDamageFactor) + 1))
			sendln(conn, "EEEEERRRRRROOOOOOOMMMMMmmmmm!!")
			send(conn, "Captain, the engines suffered ")
			send(conn, strconv.Itoa(damage))
			sendln(conn, " units of damage.")
			objUpdateEngineDam(username, damage, objectsdb)
		} else {
			sendln(conn, "Captain, our engines are overheating!")
		}
	} else {
		if rand.Float32() <= .05 {
			//this is puking
			damage := rand.Intn(int((actualMoves * WarpEngineDamageFactor) + 1))
			sendln(conn, "EEEEERRRRRROOOOOOOMMMMMmmmmm!!")
			send(conn, "Captain, the engines suffered ")
			send(conn, strconv.Itoa(damage))
			sendln(conn, " units of damage.")
			objUpdateEngineDam(username, damage, objectsdb)
		}
	}
	//
	doGameDelay(username, nummoves, objectsdb)
	return false
}

//
// Do the default Move command
//
func dodefMove(invalidcommand bool, conn *Conn) bool {
	invalidcommand = processHelp("help move", invalidcommand, conn)
	return false
}

//
// News command
//
func processNews(comnd string, invalidcommand bool, conn *Conn) bool {
	sendln(conn, "Error: The news command takes no parameters")
	return false
}

//
// Do the default News command
//
func dodefNews(invalidcommand bool, conn *Conn) bool {
	// Open file for reading
	file, err := os.Open("/home/newmanh/go/src/decwars/decwars.nws")
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	s := string(data[:])
	sendln(conn, s)
	return false
}

//
// Phasers command
// command variants:
// ph locx locy
// ph re/ab locx locy
// ph computed enemy/friendly.../name
// ph 500 locx locy
// ph 500 re/ab locx locy
// ph 500 computed enemy/friendly.../name
//
// Phasers must be directed at a specific target, and only one target may be specified per command. Obstacles seemingly
// in the path of the phaser blast are unaffected, since the energy ray isnot a line-of- sight weapon. The size of the
// hit is inversely proportional to the distance from the target. Maximum range is 10 sectors vertically, horizontally,
// or diagonally. Each phaserblast consumes 200 units of ship energy, unless a specificamount of energy is given (the
// specified energy must be between50 and 500 units, inclusive). The phaser banks have roughly a 5% chance of damage
// with a default (200 unit) blast, with the probability of damage reaching nearly 65% with a maximum (500unit) blast.
// The severity of the resulting damage is also dependant on the size of the blast.  Also, if your ship's shields are up,
// a high-speed shield control is used to quickly lower and then restore the shields during the fire.  This procedure
// consumes another 200 units of ship energy. The weapons officer on board your ship will cancel all phaser blasts
// directed against friendly ships, bases, or planets. Firing phasers (or getting hit by phasers) puts you on red alert.
// NOTE: Although phasers can damage enemy planetary installations (BUILDs), they can NOT destroy the planet itself.
//
func processPhasers(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string, pointsdb *sql.DB, usersdb *sql.DB) bool {
	var myLocx int
	var myLocy int
	var myShldUp int
	var myShipEnergy int
	var myIOFmt int
	var mySide int
	var myPhasDam int
	var myCmpDam int
	var myOutputLen int
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var parameter4 string
	var locx int
	var locy int
	var whoNme string
	var whoLocx int
	var whoLocy int
	var whoShldUp int
	var whoShipEnergy int
	var whoType int
	var whoSide int
	var err error
	var err1 error
	var rslt bool
	var pival int
	var phasamt int
	var highSpeedActivated bool
	var whoShld int
	var whoWarpEngDam int
	var whoImpEngDam int
	var whoPhoTorDam int
	var whoPhasDam int
	var whoShldDam int
	var whoCmpDam int
	var whoLifeSupDam int
	var whoRadioDam int
	var whoTractorDam int
	var whoShipDam int
	var whoObjtype int

	// Default amount to shoot
	phasamt = defPhasAmt
	//
	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	// First parm
	parameter1 = strings.Trim(prm1[1], ctlsp)
	// Lookup me
	objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy, IOFmt, Side, PhasDam, CmpDam, OutputLen from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myShipEnergy, &myIOFmt, &mySide, &myPhasDam, &myCmpDam, &myOutputLen)
	if len(prm1) == 3 { // 2 parm - ph x y
		parameter2 = strings.Trim(prm1[2], ctlsp)

		// Lookup user
		locx, err = strconv.Atoi(parameter1)
		locy, err1 = strconv.Atoi(parameter2)

		// if first parm is alpha, then 2nd must be alpha (for ph com side/name)
		if err != nil && err1 != nil {
			if testCommandMatch(NameComputed, parameter1, lenOfString1) {
				locx, locy, rslt = docomputed(username, parameter2, objectsdb)
				if rslt {
					parameter1 = strconv.Itoa(locx)
					parameter2 = strconv.Itoa(locy)
					myIOFmt = IOFmtAbs
				}
			}
		} else { //must be numeric for both values

			if myIOFmt == IOFmtAbs {
				err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
			} else {
				err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
			}
			if err != nil { // no such obj
				sendln(conn, noTarg)
				return false
			} else {
				locx = whoLocx
				locy = whoLocy
			}
		}
	} else {
		if len(prm1) == 4 { // 3 parm - ph re/ab x y OR ph 500 x y
			phasamt, err = strconv.Atoi(parameter1)
			if err == nil { // got an phasamt
				parameter2 = strings.Trim(prm1[2], ctlsp)
				parameter3 = strings.Trim(prm1[3], ctlsp)

				// Lookup user
				locx, err = strconv.Atoi(parameter2)
				locy, err1 = strconv.Atoi(parameter3)
				// if 2nd parm is alpha, then 3rd must be alpha (for ph com side/name)
				if err != nil && err1 != nil {
					if testCommandMatch(NameComputed, parameter2, lenOfString1) {
						locx, locy, rslt = docomputed(username, parameter3, objectsdb)
						if rslt {
							parameter1 = strconv.Itoa(locx)
							parameter2 = strconv.Itoa(locy)
							myIOFmt = IOFmtAbs
						}
					}
				} else { //must be numeric for both values

					if myIOFmt == IOFmtAbs {
						err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
					} else {
						err = objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, ShipEnergy,Objtype from objects WHERE Locx = ? and Locy = ?", myLocx+locx, myLocy+locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShldUp, &whoShipEnergy, &whoType)
					}
					if err != nil { // no such obj
						sendln(conn, noTarg)
						return false
					} else {
						locx = whoLocx
						locy = whoLocy

					}
				}

			} else {
				phasamt = defPhasAmt
				if testCommandMatch(NameRelative, parameter1, lenOfString1) {
					parameter2 = strings.Trim(prm1[2], ctlsp)
					parameter3 = strings.Trim(prm1[3], ctlsp)
					pival, err = strconv.Atoi(parameter2)
					locx = myLocx + pival
					pival, err1 = strconv.Atoi(parameter3)
					locy = myLocy + pival
					if err != nil || err1 != nil {
						sendln(conn, invparm)
						return false
					}
				} else {
					if testCommandMatch(NameAbsolute, parameter1, lenOfString1) {
						parameter2 = strings.Trim(prm1[2], ctlsp)
						parameter3 = strings.Trim(prm1[3], ctlsp)
						pival, err = strconv.Atoi(parameter2)
						locx = pival
						pival, err1 = strconv.Atoi(parameter3)
						locy = pival
						if err != nil || err1 != nil {
							sendln(conn, invparm)
							return false
						}
					}
				}
			}
		} else {
			if len(prm1) == 5 { // 4 parm - ph 500 re/ab x y
				phasamt, err = strconv.Atoi(parameter1)
				if err == nil { // got an phasamt
					parameter2 = strings.Trim(prm1[2], ctlsp)
					parameter3 = strings.Trim(prm1[3], ctlsp)
					parameter4 = strings.Trim(prm1[4], ctlsp)
					if testCommandMatch(NameRelative, parameter2, lenOfString1) {
						pival, err = strconv.Atoi(parameter3)
						locx = myLocx + pival
						pival, err1 = strconv.Atoi(parameter4)
						locy = myLocy + pival
						if err != nil || err1 != nil {
							sendln(conn, invparm)
							return false
						}
					} else {
						if testCommandMatch(NameAbsolute, parameter2, lenOfString1) {
							pival, err = strconv.Atoi(parameter3)
							locx = pival
							pival, err1 = strconv.Atoi(parameter4)
							locy = pival
							if err != nil || err1 != nil {
								sendln(conn, invparm)
								return false
							}
						}
					}

				}
			} else {
				sendln(conn, invparm)
				return false
			}
		}
	}

	//
	// Damage to phaser control test here
	//
	if myPhasDam > MaxDam {
		sendln(conn, phasDamMsg)
		return false
	}

	//
	// Damage to computer test here
	//
	if myCmpDam > MaxDam {
		sendln(conn, cmpDamMsg)
		return false
	}

	//
	// Get the name for the locs
	//
	err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, Side from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShld, &whoShldUp, &whoWarpEngDam, &whoImpEngDam, &whoPhoTorDam, &whoPhasDam, &whoShldDam, &whoCmpDam, &whoLifeSupDam, &whoRadioDam, &whoTractorDam, &whoShipDam, &whoObjtype, &whoSide)
	if err != nil {
		sendln(conn, noTarg)
		return false
	}

	// Is phaser amount > 50 or < 500 then error
	if phasamt < 50 || phasamt > 500 {
		sendln(conn, phasSize)
		return false
	}

	// The weapons officer on board your ship will cancel all phaser blasts
	// directed against friendly ships, bases, or planets
	if whoSide == mySide {
		if myOutputLen == OutLenLong {
			send(conn, phasFriendL)
			msg := nmeaddr(whoNme, myLocx, myLocy, myIOFmt, myOutputLen, objectsdb)
			send(conn, msg)
			sendln(conn, ", Captain.")
			return false
		} else {
			send(conn, "P")
			send(conn, phasFriendM)
			msg := nmeaddr(whoNme, myLocx, myLocy, myIOFmt, myOutputLen, objectsdb)
			sendln(conn, msg)
		}
		return false
	}

	// Distance must be less than MaxScanRng sectors (10 in decwars)
	if (myLocx-locx) > MaxScanRng || (myLocy-locy) > MaxScanRng {
		sendln(conn, phasUnable)
		return false
	}
	//
	// If your shields are up, acivate high speed shield control
	//
	if myShldUp == On {
		highSpeedActivated = true
		sendln(conn, highSpeedShieldControl)
		objUpdateShldUp(Off, username, objectsdb, conn)
		objUpdateShipEnergy(myShipEnergy-100, username, objectsdb)
	} else {
		highSpeedActivated = false
	}

	// Damage to the phaser control or to the ship's computer reduces the  strength of the phaser hit
	if myPhasDam > 0 {
		phasamt = phasamt - (phasamt * (myPhasDam / MaxDam))
	}
	if myCmpDam > 0 {
		phasamt = phasamt - (phasamt * (myCmpDam / MaxDam))
	}
	if phasamt < 0 {
		phasamt = 0
	}
	//
	// Do the phaser on the object.  Damage to this phaser control or to the ship's computer reduces the strength of the phaser hit.
	// step 1: if shields are 100% they take the full amount * rnd of the hit
	// 			if shields are < 100%, the hit is shld% * rnd and distributed over all the devices
	// step 2: reduce energy by the amount of the phaser request

	doHit(whoNme, whoSide, whoLocx, whoLocy, whoShld, whoShldUp, whoWarpEngDam, whoImpEngDam, whoPhoTorDam, whoPhasDam, whoShldDam, whoCmpDam, whoLifeSupDam, whoRadioDam, whoTractorDam, whoShipDam, whoObjtype, phasamt, objectsdb, conn, phasHit, username, myLocx, myLocy, mySide, pointsdb, usersdb)
	myShipEnergy = myShipEnergy - phasamt
	objectsdb.Exec("UPDATE objects set ShipEnergy = ?  WHERE Nme = ?", myShipEnergy, username)

	//
	// Tell everyone in the vicinity of the attack tell the full message!!!!
	// xxyyzz test - this notify was not needed...
	//	notify(username, conn, myLocx, myLocy, msgPhas, whoNme, locx, locy, objectsdb, phasamt, 0)

	//
	// If high speed shield control was activated, raise shields (no need to check tractor beam here)
	//
	if highSpeedActivated == true {
		objUpdateShldUp(On, username, objectsdb, conn)
		objUpdateShipEnergy(myShipEnergy-200, username, objectsdb)
	}

	//	Handle damages TO phasers
	//  5% chance of damage with a default (200 unit) blast, with the probability of damage reaching nearly 65%
	//  with a maximum (500unit) blast.
	//  %*100* 5 + 175=blast. So (blast - 175) / 5 = probability of dam/100
	probDam := ((float32(phasamt) - 175.0) / 5.0) / 100.0
	// get a rand fraction, if it is <= probDam then damage occured. Use random # * phaser amount to determine damage
	aRndNum := rand.Float32()
	if aRndNum <= probDam {
		sendln(conn, phasOverheating1)
		sendln(conn, phasOverheating2)
		sendln(conn, phasOverheating3)
		damage := aRndNum * float32(phasamt)
		objUpdatePhasDam(username, int(damage), objectsdb)
	}
	//
	// Finally, penelize by making user wait
	// Each phaser shot takes 2 * 500 * Milliseconds per 200 units shot to  cool off.
	//
	doGameDelay(username, phasamt/200*2, objectsdb)
	return false
}

//
// Do the default Phasers command with no parms, which is help
//
func dodefPhasers(invalidcommand bool, conn *Conn) bool {
	invalidcommand = processHelp("help phas", invalidcommand, conn)
	return false
}

//
// Planets command
//
func processPlanets(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	var mySide int
	var myLocx int
	var myLocy int
	var myIOFmt int
	var myOutputLen int
	var parameter1 string
	var qry string
	var locx int
	var locy int
	var SeenbyEnemy int
	var Nme string
	// get my info
	objectsdb.QueryRow("select Side, Locx, Locy, IOFmt, OutputLen from objects where Nme=?", username).Scan(&mySide, &myLocx, &myLocy, &myIOFmt, &myOutputLen)

	prm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")
	//	prm1 := strings.SplitAfter(prm, " ")
	if len(prm) == 3 { // Must be a default address (choose rel/ab from user default)
		pival, err := strconv.Atoi(prm[1])

		if err == nil {
			locx = pival
			pival, err = strconv.Atoi(prm[2])
			if err == nil {
				locy = pival
				objectsdb.QueryRow("select IOFmt, Locx, Locy, Side from objects where Nme=?", username).Scan(&myIOFmt, &myLocx, &myLocy, &mySide)
			}
			if myIOFmt == IOFmtRel || myIOFmt == IOFmtBoth {
				if locx > MaxScanRng || locy > MaxScanRng {
					sendln(conn, OOR)
					return false
				}
				err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", myLocx+locx, myLocy+locy).Scan(&Nme, &SeenbyEnemy)
				if err == nil {
					istrue := SeenbyEnemy & mySide
					if istrue > 0 {
						qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(myLocx+locx) + " and Locy = " + strconv.Itoa(myLocy+locy) + " and Objtype = " + strconv.Itoa(TypePlanet)
						doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
					} else {
						sendln(conn, OOR)
						return false
					}
				} else {
					sendln(conn, invparm)
					return false
				}
			} else { // absolute
				if (Abs(myLocx)-Abs(locx) > MaxScanRng) || (Abs(myLocy)-Abs(locy) > MaxScanRng) {
					sendln(conn, OOR)
					return false
				}
				err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", locx, locy).Scan(&Nme, &SeenbyEnemy)
				if err == nil {
					istrue := SeenbyEnemy & mySide
					if istrue > 0 {
						qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(locx) + " and Locy = " + strconv.Itoa(locy) + " and Objtype = " + strconv.Itoa(TypePlanet)
						doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
					} else {
						sendln(conn, OOR)
						return false
					}
				}
			}
		}
	} else {
		if len(prm) == 2 { // 1 parms - can be any keyword
			parameter1 = strings.Trim(prm[1], ctlsp)
			switch true {

			case testCommandMatch(NameEnemy, parameter1, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Side != " + strconv.Itoa(mySide) + " and Objtype = " + strconv.Itoa(TypePlanet)
				doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)

			case testCommandMatch(NameSummary, parameter1, lenOfString1):
				sendln(conn, " ")
				doSummaryPlanets(conn, objectsdb)

			case testCommandMatch(NameAll, parameter1, lenOfString1):
				qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Objtype = " + strconv.Itoa(TypePlanet)
				doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)

			case testCommandMatch(NameClosest, parameter1, lenOfString1):
				clLocx, clLocy, rslt := findClosest(myLocx, myLocy, TypePlanet, 0, objectsdb)
				if rslt == true {
					qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(clLocx) + " and Locy = " + strconv.Itoa(clLocy) + " and Objtype = " + strconv.Itoa(TypePlanet)
					doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
				} else {
					sendln(conn, "No closest planet!")
					return true
				}
			}
		} else {
			if len(prm) == 4 { // 3 parms - can be any keyword
				parameter1 = strings.Trim(prm[1], ctlsp)
				switch true {
				case testCommandMatch(NameRelative, parameter1, lenOfString1):
					parameter2 := strings.Trim(prm[2], ctlsp)
					parameter3 := strings.Trim(prm[3], ctlsp)
					pival, err := strconv.Atoi(parameter2)
					if err == nil {
						locx := pival
						pival, err := strconv.Atoi(parameter3)
						if err == nil {
							locy := pival
							if locx > MaxScanRng || locy > MaxScanRng {
								sendln(conn, OOR)
								return false
							}
							qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(locx+myLocx) + " and Locy = " + strconv.Itoa(locy+myLocy) + " and Objtype = " + strconv.Itoa(TypePlanet)
							doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
						}
						return false
					}

				case testCommandMatch(NameAbsolute, parameter1, lenOfString1):
					parameter2 := strings.Trim(prm[2], ctlsp)
					parameter3 := strings.Trim(prm[3], ctlsp)
					pival, err := strconv.Atoi(parameter2)
					if err == nil {
						locx := pival
						pival, err := strconv.Atoi(parameter3)
						if err == nil {
							locy := pival
							if (myLocx-locx) > MaxScanRng || (myLocy-locy) > MaxScanRng {
								sendln(conn, OOR)
								return false
							}
							qry = "select Nme, Locx, Locy, Side, SeenbyEnemy, Objtype, Shld from objects where Locx = " + strconv.Itoa(locx) + " and Locy = " + strconv.Itoa(locy) + " and Objtype = " + strconv.Itoa(TypePlanet)
							doListQuery(invalidcommand, conn, objectsdb, username, qry, mySide, myLocx, myLocy, myIOFmt, myOutputLen)
						}
						return false
					}

				default:
					sendln(conn, invparm)
					return true

				}
				sendln(conn, invparm)
				return false
			}
		}
		//					sendln(conn, invparm)
		return false
	}
	return false

}

//
// Do the default Planets command - list all planets seen by a side
//
func dodefPlanets(username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	invalidcommand = processPlanets("pl all", invalidcommand, conn, objectsdb, username)
	doSummaryPlanets(conn, objectsdb)
	sendln(conn, " ")
	return false
}

//
// Points command for Friendly or Enemy
//
func processPoints(username string, comnd string, invalidcommand bool, pointsdb *sql.DB, objectsdb *sql.DB, usersdb *sql.DB, conn *Conn) bool {
	var mySide int

	parm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")
	//
	// Get my side so we can figure out which form of points to display
	//
	_ = objectsdb.QueryRow("select Side from objects WHERE Nme = ?", username).Scan(&mySide)

	//
	//
	//
	if testCommandMatch(NameFriendly, parm[1], lenOfString1) == true {
		if mySide == SideCoalition {
			invalidcommand = dodefPoints(username, invalidcommand, 1, conn, pointsdb, objectsdb, usersdb)
		} else {
			invalidcommand = dodefPoints(username, invalidcommand, 2, conn, pointsdb, objectsdb, usersdb)
		}
	} else {
		if testCommandMatch(NameEnemy, parm[1], lenOfString1) == true {
			if mySide == SideCoalition {
				invalidcommand = dodefPoints(username, invalidcommand, 3, conn, pointsdb, objectsdb, usersdb)
			} else {
				invalidcommand = dodefPoints(username, invalidcommand, 4, conn, pointsdb, objectsdb, usersdb)
			}
		} else {
			sendln(conn, invparm)
			return false
		}
	}
	return false

}

//
// Do the default Points command
//                         You     Coalition       Empire  Archerons       Neutrals        Career Totals
// Damage to enemy         0       2               2       6               1               0
// Enemy destroyed         0       2               2       6               1               0
// Damage to bases         0       2               2       6               1               0
// Bases destroyed         0       2               2       6               1               0
// Damage to archerons     0       2               2       6               1               0
// Archerons destroyed     0       2               2       6               1               0
// Damage to stars         0       2               2       6               1               0
// Stars destroyed         0       2               2       6               1               0
// Damage to planets       0       2               2       6               1               0
// Planets destroyed       0       2               2       6               1               0
//
// Total Points            0       2               2       6               1               0
//
// Number of ships         0       2               2       6               1               0
// Number stardates:       1       2               2       6               1               0
// Points per player       0       2               2       6               1               0
// Points per stardate     0       2               2       6               1               0
//
//
// Congratulations, Captain! You
// are now in first place!
// Display: 0 = all
//          1 = Coalition alone (friendly, side=coalition)
//          2 = Empire alone (friendly, side=Empire)
//          3 = Coalition + robots  (enemy, side=Empire)
//          4 = Empire + robots (enemy, side=coalition)

func dodefPoints(username string, invalidcommand bool, disp int, conn *Conn, pointsdb *sql.DB, objectsdb *sql.DB, usersdb *sql.DB) bool {
	var myDamToBases int
	var myDamToBH int
	var myDamToShip int
	var myDamToStars int
	var myDamToPlanets int
	var myNumOfShips int
	var myNumOfStarDates int

	var StarDate int

	//
	// Points stuff for usersdb
	//
	var uDamToBases int
	var uDamToBH int
	var uDamToShip int
	var uDamToStars int
	var uDamToPlanets int
	var uNumOfShips int
	var uNumOfStarDates int

	// for Totals
	var cDamToBH int
	var cDamToBases int
	var cDamToShip int
	var cDamToStars int
	var cDamToPlanets int
	var cNumOfShips int
	var cNumOfStarDates int

	var eDamToBases int
	var eDamToBH int
	var eDamToShip int
	var eDamToStars int
	var eDamToPlanets int
	var eNumOfShips int
	var eNumOfStarDates int

	var aDamToBases int
	var aDamToBH int
	var aDamToShip int
	var aDamToStars int
	var aDamToPlanets int
	var aNumOfShips int
	var aNumOfStarDates int

	var nDamToBases int
	var nDamToBH int
	var nDamToShip int
	var nDamToStars int
	var nDamToPlanets int
	var nNumOfShips int
	var nNumOfStarDates int

	var nCumWinE int
	var nCumWinC int

	// Get info required for each line
	pointsdb.QueryRow("select DamToBH, DamToBases, DamToShip, DamToStars, DamToPlanets, NumOfShips, NumOfStarDates  from points where Nme=?", username).Scan(&myDamToBH, &myDamToBases, &myDamToShip, &myDamToStars, &myDamToPlanets, &myNumOfShips, &myNumOfStarDates)

	objectsdb.QueryRow("select StarDate from objects where Nme = ?", username).Scan(&StarDate)

	usersdb.QueryRow("select DamToBH, DamToBases, DamToShip, DamToStars, DamToPlanets, NumOfShips, NumOfStarDates from users where name=?", username).Scan(&uDamToBH, &uDamToBases, &uDamToShip, &uDamToStars, &uDamToPlanets, &uNumOfShips, &uNumOfStarDates)
	//
	// Points have lots of summary amounts, do it here
	//
	pointsdb.QueryRow("select total(DamToBH), total(DamToBases), total(DamToShip), total(DamToStars), total(DamToPlanets), total(NumOfShips), total(NumOfStarDates)  from points where Side=?", SideCoalition).Scan(&cDamToBH, &cDamToBases, &cDamToShip, &cDamToStars, &cDamToPlanets, &cNumOfShips, &cNumOfStarDates)

	pointsdb.QueryRow("select total(DamToBH), total(DamToBases),  total(DamToShip), total(DamToStars), total(DamToPlanets), total(NumOfShips), total(NumOfStarDates)  from points where Side=?", SideEmpire).Scan(&eDamToBH, &eDamToBases, &eDamToShip, &eDamToStars, &eDamToPlanets, &eNumOfShips, &eNumOfStarDates)

	pointsdb.QueryRow("select total(DamToBH), total(DamToBases), total(DamToShip), total(DamToStars), total(DamToPlanets), total(NumOfShips), total(NumOfStarDates)  from points where Side=?", SideArcheron).Scan(&aDamToBH, &aDamToBases, &aDamToShip, &aDamToStars, &aDamToPlanets, &aNumOfShips, &aNumOfStarDates)

	pointsdb.QueryRow("select total(DamToBH), total(DamToBases), total(DamToShip), total(DamToStars), total(DamToPlanets), total(NumOfShips), total(NumOfStarDates)  from points where Side=?", SideNeutral).Scan(&nDamToBH, &nDamToBases, &nDamToShip, &nDamToStars, &nDamToPlanets, &nNumOfShips, &nNumOfStarDates)

	usersdb.QueryRow(`select CumlativeWins from users where name = "empire";`).Scan(&nCumWinE)
	usersdb.QueryRow(`select CumlativeWins from users where name = "coalit";`).Scan(&nCumWinC)
	// Header
	if disp == 0 {
		sendln(conn, "\t\t\tYou\tCoalition\tEmpire\t\tArcherons\tNeutrals\tCareer Totals")
	} else {
		if disp == 1 {
			sendln(conn, "\t\t\tYou\tCoalition\tCareer Totals")
		} else {
			if disp == 2 {
				sendln(conn, "\t\t\tYou\tEmpire\t\tCareer Totals")
			} else {
				if disp == 3 {
					sendln(conn, "\t\t\tYou\tEmpire\t\tArcherons\tNeutrals\tCareer Totals")
				} else {
					sendln(conn, "\t\t\tYou\tCoalition\tArcherons\tNeutrals\tCareer Totals")
				}
			}
		}
	}
	send(conn, "Damage to ships:\t")
	send(conn, strconv.Itoa(myDamToShip))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cDamToShip))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eDamToShip))
		send(conn, "\t\t")

		send(conn, strconv.Itoa(aDamToShip))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nDamToShip))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(uDamToShip))

	send(conn, "Damage to planets:\t")
	send(conn, strconv.Itoa(myDamToPlanets))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cDamToPlanets))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eDamToPlanets))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(aDamToPlanets))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nDamToPlanets))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(uDamToPlanets))

	send(conn, "Damage to stars:\t")
	send(conn, strconv.Itoa(myDamToStars))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cDamToStars))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eDamToStars))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(aDamToStars))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nDamToStars))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(uDamToStars))

	send(conn, "Damage to bases:\t")
	send(conn, strconv.Itoa(myDamToBases))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cDamToBases))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eDamToBases))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(aDamToBases))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nDamToBases))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(uDamToBases))

	send(conn, "Damage to black holes:\t")
	send(conn, strconv.Itoa(myDamToBH))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cDamToBH))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eDamToBH))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(aDamToBH))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nDamToBH))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(uDamToBH))

	sendln(conn, " ")
	send(conn, "Total Points:\t\t")
	send(conn, strconv.Itoa(myDamToBases+myDamToShip+myDamToStars+myDamToPlanets+myDamToBH))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cDamToBases+cDamToShip+cDamToStars+cDamToPlanets+cDamToBH))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eDamToBases+eDamToShip+eDamToStars+eDamToPlanets+eDamToBH))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(aDamToBases+aDamToShip+aDamToStars+aDamToPlanets+aDamToBH))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nDamToBases+nDamToShip+nDamToStars+nDamToPlanets+nDamToBH))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(uDamToBases+uDamToShip+uDamToStars+uDamToPlanets+uDamToBH))

	sendln(conn, " ")
	send(conn, "Number of targets:\t")
	send(conn, strconv.Itoa(myNumOfShips))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cNumOfShips))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eNumOfShips))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(aNumOfShips))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nNumOfShips))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(uNumOfShips))

	send(conn, "Number stardates:\t")
	send(conn, strconv.Itoa(StarDate))
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		send(conn, strconv.Itoa(cNumOfStarDates))
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		send(conn, strconv.Itoa(eNumOfStarDates))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(aNumOfStarDates))
		send(conn, "\t\t")
		send(conn, strconv.Itoa(nNumOfStarDates))
		send(conn, "\t\t")
	}
	sendln(conn, strconv.Itoa(myNumOfStarDates))

	send(conn, "Points per target:\t")
	if myNumOfShips != 0 {
		send(conn, strconv.Itoa((myDamToBases+myDamToShip+myDamToStars+myDamToPlanets+myDamToBH)/myNumOfShips))
	} else {
		send(conn, strconv.Itoa(0))
	}
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		if cNumOfShips != 0 {
			send(conn, strconv.Itoa((cDamToBases+cDamToShip+cDamToStars+cDamToPlanets+cDamToBH)/cNumOfShips))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		if eNumOfShips != 0 {
			send(conn, strconv.Itoa((eDamToBases+eDamToShip+eDamToStars+eDamToPlanets+eDamToBH)/eNumOfShips))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
		if aNumOfShips != 0 {
			send(conn, strconv.Itoa((aDamToBases+aDamToShip+aDamToStars+aDamToPlanets+aDamToBH)/aNumOfShips))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
		if nNumOfShips != 0 {
			send(conn, strconv.Itoa((nDamToBases+nDamToShip+nDamToStars+nDamToPlanets+nDamToBH)/nNumOfShips))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
	}
	if uNumOfShips != 0 {
		sendln(conn, strconv.Itoa((uDamToBases+uDamToShip+uDamToStars+uDamToPlanets+uDamToBH)/uNumOfShips))
	} else {
		sendln(conn, strconv.Itoa(0))
	}

	send(conn, "Points per stardate:\t")
	if StarDate != 0 {
		send(conn, strconv.Itoa((myDamToBases+myDamToShip+myDamToStars+myDamToPlanets+myDamToBH)/StarDate))
	} else {
		send(conn, strconv.Itoa(0))
	}
	send(conn, "\t")
	if disp == 0 || disp == 1 || disp == 2 {
		if cNumOfStarDates != 0 {
			send(conn, strconv.Itoa((cDamToBases+cDamToShip+cDamToStars+cDamToPlanets+cDamToBH)/cNumOfStarDates))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
	}
	if disp == 0 || disp == 3 || disp == 4 {
		if eNumOfStarDates != 0 {
			send(conn, strconv.Itoa((eDamToBases+eDamToShip+eDamToStars+eDamToPlanets+eDamToBH)/eNumOfStarDates))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
		if aNumOfStarDates != 0 {
			send(conn, strconv.Itoa((aDamToBases+aDamToShip+aDamToStars+aDamToPlanets+aDamToBH)/aNumOfStarDates))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
		if nNumOfStarDates != 0 {
			send(conn, strconv.Itoa((nDamToBases+nDamToShip+nDamToStars+nDamToPlanets+nDamToBH)/nNumOfStarDates))
		} else {
			send(conn, strconv.Itoa(0))
		}
		send(conn, "\t\t")
	}
	if myNumOfStarDates != 0 {
		sendln(conn, strconv.Itoa((uDamToBases+uDamToShip+uDamToStars+uDamToPlanets+uDamToBH)/myNumOfStarDates))
	} else {
		sendln(conn, strconv.Itoa(0))
	}

	send(conn, "Cummulative Wins:\t\t")
	send(conn, strconv.Itoa(nCumWinC))
	send(conn, "\t\t")
	sendln(conn, strconv.Itoa(nCumWinE))
	return false
}

//
// Radio command - with or without parms - radio on/off
//
func processRadio(comnd string, username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var myRadio int
	var parameter1 string

	// get my info
	objectsdb.QueryRow("select Radio from objects where Nme=?", username).Scan(&myRadio)

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")

	// First parm
	parameter1 = strings.Trim(prm1[1], ctlsp)

	if len(prm1) == 2 {
		switch true {
		case testCommandMatch(NameOn, parameter1, lenOfString2):
			if myRadio == 1 {
				sendln(conn, "Radio is already on, captain")
			} else {
				myRadio = 1
				objectsdb.Exec("UPDATE objects set Radio = ? WHERE Nme = ?", myRadio, username)
			}

		case testCommandMatch(NameOff, parameter1, lenOfString2):
			if myRadio == 0 {
				sendln(conn, "Radio is already off, captain")
			} else {
				myRadio = 0
				objectsdb.Exec("UPDATE objects set Radio = ? WHERE Nme = ?", myRadio, username)
			}

		default:
			sendln(conn, invparm)
		}
	} else {
		sendln(conn, invparm)
	}

	return false
}

//
// Radio command - Without parms = display setting
//
func dodefRadio(comnd string, username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var myRadio int
	// get my info
	objectsdb.QueryRow("select Radio from objects where Nme=?", username).Scan(&myRadio)
	doStatusRadio(conn, myRadio)
	return false
}

//
// Do a repair needs if docked/not docked rates
//
func doRepair(username string, repairAmt int, objectsdb *sql.DB) {
	var WarpEngDam int
	var ImpEngDam int
	var PhoTorDam int
	var PhasDam int
	var ShldDam int
	var CmpDam int
	var LifeSupDam int
	var RadioDam int
	var TractorDam int
	var Objtype int
	var DockFlag string
	var Side int

	//
	// Step 1: get the current status of the user being hit
	//
	objectsdb.QueryRow("select WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, Objtype, DockFlag, Side from objects WHERE Nme = ?", username).Scan(&WarpEngDam, &ImpEngDam, &PhoTorDam, &PhasDam, &ShldDam, &CmpDam, &LifeSupDam, &RadioDam, &TractorDam, &Objtype, &DockFlag, &Side)

	//
	// Step 2: If the damage > repairAmt, reduce by repairAmt, else dam = 0
	//
	if WarpEngDam > repairAmt {
		WarpEngDam = WarpEngDam - repairAmt
	} else {
		WarpEngDam = 0
	}
	if ImpEngDam > repairAmt {
		ImpEngDam = ImpEngDam - repairAmt
	} else {
		ImpEngDam = 0
	}
	if PhoTorDam > repairAmt {
		PhoTorDam = PhoTorDam - repairAmt
	} else {
		PhoTorDam = 0
	}
	if PhasDam > repairAmt {
		PhasDam = PhasDam - repairAmt
	} else {
		PhasDam = 0
	}
	if ShldDam > repairAmt {
		ShldDam = ShldDam - repairAmt
	} else {
		ShldDam = 0
	}
	if CmpDam > repairAmt {
		CmpDam = CmpDam - repairAmt
	} else {
		CmpDam = 0
	}

	// should only fix life support if docked
	if DockFlag != "" {
		if LifeSupDam > repairAmt {
			LifeSupDam = LifeSupDam - repairAmt
		} else {
			LifeSupDam = 0
		}
	}

	if RadioDam > repairAmt {
		RadioDam = RadioDam - repairAmt
	} else {
		RadioDam = 0
	}
	if TractorDam > repairAmt {
		TractorDam = TractorDam - repairAmt
	} else {
		TractorDam = 0
	}
	//
	// Step 3: Update the current status of the user being hit
	//
	objectsdb.Exec("UPDATE objects set WarpEngDam=?, ImpEngDam=?, PhoTorDam=?, PhasDam=?, ShldDam=?, CmpDam=?, LifeSupDam=?, RadioDam=?, TractorDam=?  WHERE Nme = ?", WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, username)
	// do not delay if it's a Ancheron
	if Side != SideArcheron {
		time.Sleep(time.Duration(1) * time.Duration(movedelay) * time.Millisecond)
	}
}

//
// Repair command - options: [radio|shields|tractor|Warp|Impulse|Photon|Phaser|Computer|Life] - use all repair amount on 1 device
//
func processRepair(username string, comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var WarpEngDam int
	var ImpEngDam int
	var PhoTorDam int
	var PhasDam int
	var ShldDam int
	var CmpDam int
	var LifeSupDam int
	var RadioDam int
	var TractorDam int
	var Objtype int
	var DockFlag string

	repairAmt := 50
	sndto := strings.Split(strings.TrimSpace(strings.Trim(comnd, ctlonly)), " ")
	//
	// Step 1: get the current status of the user being hit
	//
	objectsdb.QueryRow("select WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, Objtype, DockFlag from objects WHERE Nme = ?", username).Scan(&WarpEngDam, &ImpEngDam, &PhoTorDam, &PhasDam, &ShldDam, &CmpDam, &LifeSupDam, &RadioDam, &TractorDam, &Objtype, &DockFlag)

	// just asked for help if invalid parms
	if len(sndto) != 2 {
		invalidcommand = processHelp("help repair", invalidcommand, conn)
		return false
	} else {
		switch true {
		case testCommandMatch("radio", sndto[1], lenOfString2):
			if RadioDam > repairAmt {
				RadioDam = RadioDam - repairAmt
			} else {
				RadioDam = 0
			}
		case testCommandMatch("shields", sndto[1], lenOfString2):
			if ShldDam > repairAmt {
				ShldDam = ShldDam - repairAmt
			} else {
				ShldDam = 0
			}
		case testCommandMatch("warp", sndto[1], lenOfString2):
			if WarpEngDam > repairAmt {
				WarpEngDam = WarpEngDam - repairAmt
			} else {
				WarpEngDam = 0
			}
		case testCommandMatch("tractor", sndto[1], lenOfString2):
			if TractorDam > repairAmt {
				TractorDam = TractorDam - repairAmt
			} else {
				TractorDam = 0
			}
		case testCommandMatch("impulse", sndto[1], lenOfString2):
			if ImpEngDam > repairAmt {
				ImpEngDam = ImpEngDam - repairAmt
			} else {
				ImpEngDam = 0
			}
		case testCommandMatch("photon", sndto[1], lenOfString3):
			if PhoTorDam > repairAmt {
				PhoTorDam = PhoTorDam - repairAmt
			} else {
				PhoTorDam = 0
			}
		case testCommandMatch("phaser", sndto[1], lenOfString3):
			if PhasDam > repairAmt {
				PhasDam = PhasDam - repairAmt
			} else {
				PhasDam = 0
			}
		case testCommandMatch("computer", sndto[1], lenOfString2):
			if CmpDam > repairAmt {
				CmpDam = CmpDam - repairAmt
			} else {
				CmpDam = 0
			}
		case testCommandMatch("life", sndto[1], lenOfString2):
			// should only fix life support if docked
			if DockFlag != "" {
				if LifeSupDam > repairAmt {
					LifeSupDam = LifeSupDam - repairAmt
				} else {
					LifeSupDam = 0
				}
			}

		default:
			sendln(conn, invparm)
			return true
		}

	}
	//
	// Step 3: Update the current status of the user being hit
	//
	objectsdb.Exec("UPDATE objects set WarpEngDam=?, ImpEngDam=?, PhoTorDam=?, PhasDam=?, ShldDam=?, CmpDam=?, LifeSupDam=?, RadioDam=?, TractorDam=?  WHERE Nme = ?", WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, username)
	time.Sleep(time.Duration(1) * time.Duration(movedelay) * time.Millisecond)
	return false
}

//
// Do the default Repair command - no parms
// 30 units per turn, 50 on a requested repair, 100 on a docked repair
// Needs to check for ship being docked xxxxxxxxxxxxxxxxxxxxxxxxxxxx
//
func dodefRepair(username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	doRepair(username, 50, objectsdb)
	//
	// finally do a delay
	//
	doGameDelay(username, 1, objectsdb)
	return false
}

//
// Calc if it is an even number if even return 1
//
func Even(number int) bool {
	return number%2 == 0
}

//
// Do the actual scan work
// Input: top left corner, top right corner, bot left corner, bot right corner, flag if to show warnings, format-rel or abs
// ****  needs to have the warnings flag work ***
//
func dotheScan(conn *Conn, username string, tlcx int, tlcy int, trcx int, trcy int, blcx int, blcy int, brcx int, brcy int, wrn bool, formt int, Locx int, Locy int, SeenbyEnemy int, mySide int, objectsdb *sql.DB) {
	var Objtype int
	var Side int
	//	var h int
	var h1 int
	//	var r int
	var r1 int

	// Print header -
	headerlow := tlcy
	headerhigh := trcy

	// Print header
	send(conn, "      ")
	if formt == IOFmtAbs {
		if Even(headerlow) == true {
			send(conn, " ")
		} else {
			send(conn, "  ")
		}
	}
	for h1 = headerlow; h1 <= headerhigh; h1++ {
		if Even(h1) == true {
			if formt == IOFmtAbs {
				send(conn, strconv.Itoa(h1))
				if h1 < -9 || h1 > 9 {
					send(conn, " ")
				} else {
					if h1 >= -10 || h1 <= 10 {
						send(conn, "  ")
					} else {
						send(conn, " ")
					}
				}
			} else {
				send(conn, strconv.Itoa(h1-Locy))
				if (h1 - Locy) < -2 {
					send(conn, " ")
				} else {
					if (h1-Locy) >= -2 && (h1-Locy) <= 1 {
						send(conn, "  ")
					} else {
						send(conn, "  ")
					}
				}
			}
			if formt == IOFmtAbs {
				send(conn, "  ")
			}
		} else {
			if formt == IOFmtAbs {
				send(conn, "   ")
			} else {
				send(conn, "     ")
			}
		}
	}
	sendln(conn, "")

	// Print 1 row at a time
	rowlow := blcx
	rowhigh := trcx

	// Set the seenbyenemy if seen
	for r1 = rowhigh; r1 >= rowlow; r1-- {
		if formt == IOFmtAbs {
			if r1 > 9 || r1 < -9 {
				send(conn, "    ")
			} else {
				send(conn, "     ")
			}
			send(conn, strconv.Itoa(r1))
		} else {
			if r1-Locx > 9 {
				send(conn, "    ")
			} else {
				if r1-Locx < -9 {
					send(conn, "   ")
				} else {
					if r1-Locx < 0 {
						send(conn, "    ")
					} else {
						send(conn, "     ")
					}
				}
			}
			send(conn, strconv.Itoa(r1-Locx))
		}
		for h1 = headerlow; h1 <= headerhigh; h1++ {
			err := objectsdb.QueryRow("select Objtype, Side, SeenbyEnemy from objects WHERE Locx = ? and Locy = ?", r1, h1).Scan(&Objtype, &Side, &SeenbyEnemy)
			if err == nil {
				//
				// Got an object - if their seenbye doesn't have my side you must add it here
				//
				//Compare if it's been seen by myside
				istrue := SeenbyEnemy & mySide
				if istrue == 0 {
					SeenbyEnemy = SeenbyEnemy | mySide
					objectsdb.Exec("UPDATE objects set SeenbyEnemy = ? WHERE Locx = ? and Locy = ?", SeenbyEnemy, r1, h1)
				}

				//
				if Objtype == TypeShip {
					if Side == SideCoalition {
						send(conn, "  +C")
					} else {
						if Side == SideEmpire {
							send(conn, "  -E")
						} else {
							if Side == SideNeutral {
								send(conn, "  NN")
							} else {
								send(conn, "  ^^")
							}
						}
					}
				} else {
					if Objtype == TypePlanet {
						if Side == SideNeutral {
							send(conn, "  @ ")
						} else {
							if Side == SideCoalition {
								send(conn, "  +@")
							} else {
								if Side == SideEmpire {
									send(conn, "  -@")
								} else {
									send(conn, "  ^@")
								}
							}
						}
					} else {
						if Objtype == TypeStar {
							send(conn, "  * ")
						} else {
							if Objtype == TypeBase {
								if Side == SideArcheron {
									send(conn, "  ^B")
								} else {
									if Side == SideCoalition {
										send(conn, "  +B")
									} else {
										if Side == SideEmpire {
											send(conn, "  -B")
										}
									}
								}
							} else {
								if Objtype == TypeBH {
									send(conn, "    ")
								}
							}
						}
					}
				}
			} else {
				if wrn == true {
					//
					// Warning for bases & planets
					//
					v1 := h1 - MaxBaseRng
					v2 := h1 + MaxBaseRng
					v3 := r1 - MaxBaseRng
					v4 := r1 + MaxBaseRng
					p1 := h1 - MaxPlanetRng
					p2 := h1 + MaxPlanetRng
					p3 := r1 - MaxPlanetRng
					p4 := r1 + MaxPlanetRng
					var cnt int
					objectsdb.QueryRow("select count(Nme) from objects where Side != ? and ((Objtype = ? and ((Locx between ? and ?) and (Locy between ? and ?)) or (Objtype = ? and ((Locx between ? and ?) and (Locy between ? and ?)))));",
						mySide, TypeBase, v3, v4, v1, v2, TypePlanet, p3, p4, p1, p2).Scan(&cnt)
					if cnt == 0 {
						send(conn, "  . ")
					} else {
						send(conn, "  ! ")
					}
				} else {
					send(conn, "  . ")
				}
			}
		}
		send(conn, "  ")
		if formt == IOFmtAbs || formt == IOFmtBoth {
			send(conn, "  ")
			sendln(conn, strconv.Itoa(r1))
		} else {
			send(conn, "  ")
			sendln(conn, strconv.Itoa(r1-Locx))
		}
	}
	// Print trailer
	send(conn, "      ")
	if formt == IOFmtAbs || formt == IOFmtBoth {
		if Even(headerlow) == true {
			send(conn, " ")
		} else {
			send(conn, "  ")
		}
	}
	for h1 = headerlow; h1 <= headerhigh; h1++ {
		if Even(h1) == true {
			if formt == IOFmtAbs || formt == IOFmtBoth {
				send(conn, strconv.Itoa(h1))
				if h1 < -9 || h1 > 9 {
					send(conn, " ")
				} else {
					if h1 >= -10 || h1 <= 10 {
						send(conn, "  ")
					} else {
						send(conn, " ")
					}
				}
			} else {
				send(conn, strconv.Itoa(h1-Locy))
				if (h1 - Locy) < -2 {
					send(conn, " ")
				} else {
					if (h1-Locy) >= -2 && (h1-Locy) <= 1 {
						send(conn, "  ")
					} else {
						send(conn, "  ")
					}
				}
			}
			if formt == IOFmtAbs || formt == IOFmtBoth {
				send(conn, "  ")
			}
		} else {
			if formt == IOFmtAbs || formt == IOFmtBoth {
				send(conn, "   ")
			} else {
				send(conn, "     ")
			}
		}
	}
	sendln(conn, "")
}

//
// Scan command this one takes parms - [Up|Down|Right|Left|Corner] [<range>|<vr><hr>] [Warning]
//
func processScan(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB) bool {
	var Nme string
	var Locx int
	var Locy int
	var IOFmt int
	var SeenbyEnemy int
	var Side int

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")

	// Get users' info for doing scan
	_ = objectsdb.QueryRow("select Nme, Locx, Locy, IOFmt, SeenbyEnemy, Side from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &IOFmt, &SeenbyEnemy, &Side)

	headerlow := Locy - MaxScanRng
	headerhigh := Locy + MaxScanRng
	if headerlow < 0 {
		headerlow = 0
	}
	if headerhigh > Vmax {
		headerhigh = Vmax
	}
	rowlow := Locx - MaxScanRng
	rowhigh := Locx + MaxScanRng
	if rowlow < 0 {
		rowlow = 0
	}
	if rowhigh > Hmax {
		rowhigh = Hmax
	}

	if len(prm1) == 2 { //1 parms
		parameter1 := strings.Trim(prm1[1], ctlsp)
		// parm can be alpha for 1 parm - up down right left or W or 1 number
		if testCommandMatch("up", parameter1, lenOfString1) {
			dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, Locx, headerlow, Locx, headerhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
		} else {
			if testCommandMatch("down", parameter1, lenOfString1) {
				dotheScan(conn, username, Locx, headerlow, Locx, headerhigh, rowlow, headerlow, rowlow, headerhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
			} else {
				if testCommandMatch("right", parameter1, lenOfString1) {
					dotheScan(conn, username, rowhigh, Locy, rowhigh, headerhigh, rowlow, Locy, rowlow, headerhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
				} else {
					if testCommandMatch("left", parameter1, lenOfString1) {
						dotheScan(conn, username, rowhigh, headerlow, rowhigh, Locy, rowlow, headerlow, rowlow, Locy, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
					} else {
						if testCommandMatch("warning", parameter1, lenOfString1) {
							dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, rowlow, headerlow, rowlow, headerhigh, true, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
						} else { //numeric?
							pival, err := strconv.Atoi(parameter1)
							if err == nil {
								if pival > MaxScanRng {
									sendln(conn, invparm)
									return false
								}
								headerlow := Locy - pival
								headerhigh := Locy + pival
								if headerlow < 0 {
									headerlow = 0
								}
								if headerhigh > Vmax {
									headerhigh = Vmax
								}
								rowlow := Locx - pival
								rowhigh := Locx + pival
								if rowlow < 0 {
									rowlow = 0
								}
								if rowhigh > Hmax {
									rowhigh = Hmax
								}
								dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, rowlow, headerlow, rowlow, headerhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
							}

						}
					}
				}
			}
		}
	} else {
		if len(prm1) == 3 { //2 parms
			// parm1 can be alpha: up down right left or corner. If parm 1 is alpha, parm 2 must be numeric or a "Warning".
			// Parm 1 can be numeric, if so parm 2 must be numeric or a W
			parameter1 := strings.Trim(prm1[1], ctlsp)
			parameter2 := strings.Trim(prm1[2], ctlsp)

			//
			// Is it numeric   **** check for max size allowed  * 2 inputs numeric
			//
			pival1, err := strconv.Atoi(parameter1)
			if err == nil {
				if pival1 > MaxScanRng {
					sendln(conn, invparm)
					return false
				}
				pival2, err := strconv.Atoi(parameter2)
				if err == nil {
					if pival2 > MaxScanRng {
						sendln(conn, invparm)
						return false
					}
					headerlow := Locx - pival1
					headerhigh := Locx + pival1
					if headerlow < 0 {
						headerlow = 0
					}
					if headerhigh > Vmax {
						headerhigh = Vmax
					}
					rowlow := Locy - pival2
					rowhigh := Locy + pival2
					if rowlow < 0 {
						rowlow = 0
					}
					if rowhigh > Hmax {
						rowhigh = Hmax
					}
					dotheScan(conn, username, headerhigh, rowlow, headerhigh, rowhigh, headerlow, rowlow, headerlow, rowhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
				} else {
					if testCommandMatch("warning", parameter2, lenOfString1) {
						headerlow := Locy - pival1
						headerhigh := Locy + pival1
						if headerlow < 0 {
							headerlow = 0
						}
						if headerhigh > Vmax {
							headerhigh = Vmax
						}
						rowlow := Locx - pival1
						rowhigh := Locx + pival1
						if rowlow < 0 {
							rowlow = 0
						}
						if rowhigh > Hmax {
							rowhigh = Hmax
						}
						dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, rowlow, headerlow, rowlow, headerhigh, true, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
					} else {
						sendln(conn, invparm)
						return true
					}
				}
			} else {
				// parm 1 is alpha
				// check parm 2 to be numeric or w
				// test for parm2=number first - if it is then setup the size for the up/down etc
				pival2, err := strconv.Atoi(parameter2)
				if err == nil {
					headerlow := Locy - pival2
					headerhigh := Locy + pival2
					if headerlow < 0 {
						headerlow = 0
					}
					if headerhigh > Vmax {
						headerhigh = Vmax
					}
					rowlow := Locx - pival2
					rowhigh := Locx + pival2
					if rowlow < 0 {
						rowlow = 0
					}
					if rowhigh > Hmax {
						rowhigh = Hmax
					}
				} else {
					// parm 2 must be a w here or error
					if !testCommandMatch("warning", parameter2, lenOfString1) {
						sendln(conn, invparm)
						return true
					}
					if testCommandMatch("up", parameter1, lenOfString1) {
						dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, Locx, headerlow, Locx, headerhigh, true, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
					} else {
						if testCommandMatch("down", parameter1, lenOfString1) {
							dotheScan(conn, username, Locx, headerlow, Locx, headerhigh, rowlow, headerlow, rowlow, headerhigh, true, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
						} else {
							if testCommandMatch("right", parameter1, lenOfString1) {
								dotheScan(conn, username, rowhigh, Locy, rowhigh, headerhigh, rowlow, Locy, rowlow, headerhigh, true, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
							} else {
								if testCommandMatch("left", parameter1, lenOfString1) {
									dotheScan(conn, username, rowhigh, headerlow, rowhigh, Locy, rowlow, headerlow, rowlow, Locy, true, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
								}
							}
						}
					}
				}
			}
		} else {
			// 3 parms?
			if len(prm1) == 4 { //3 parms
				// With 3 parms, parm 1 & 2 must be numeric and 3 must be a "Warning".

				parameter1 := strings.Trim(prm1[1], ctlsp)
				parameter2 := strings.Trim(prm1[2], ctlsp)
				parameter3 := strings.Trim(prm1[3], ctlsp)

				//
				// Is it numeric   **** check for max size allowed
				//
				pival1, err := strconv.Atoi(parameter1)
				if err == nil {
					if pival1 > MaxScanRng {
						sendln(conn, invparm)
						return false
					}
					pival2, err := strconv.Atoi(parameter2)
					if err == nil {
						if pival2 > MaxScanRng {
							sendln(conn, invparm)
							return false
						}

						headerlow := Locx - pival1
						headerhigh := Locx + pival1
						if headerlow < 0 {
							headerlow = 0
						}
						if headerhigh > Vmax {
							headerhigh = Vmax
						}
						rowlow := Locy - pival2
						rowhigh := Locy + pival2
						if rowlow < 0 {
							rowlow = 0
						}
						if rowhigh > Hmax {
							rowhigh = Hmax
						}
						if testCommandMatch("warning", parameter3, lenOfString1) {
							dotheScan(conn, username, headerhigh, rowlow, headerhigh, rowhigh, headerlow, rowlow, headerlow, rowhigh, true, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)
						} else {
							sendln(conn, invparm)
							return true
						}
					}
				}
			}
		}
	}
	return false
}

//
// Do the default Scan command - default parms - note set seenby for any object seen!!!
//
func dodefScan(invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB) bool {
	var Nme string
	var Locx int
	var Locy int
	var IOFmt int
	var SeenbyEnemy int
	var Side int

	//	doScan(conn, MaxWarpFactor+1, MaxWarpFactor+1, username, objectsdb)

	// Get users' info for doing scan
	_ = objectsdb.QueryRow("select Nme, Locx, Locy, IOFmt, SeenbyEnemy, Side from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &IOFmt, &SeenbyEnemy, &Side)

	headerlow := Locy - MaxScanRng
	headerhigh := Locy + MaxScanRng
	if headerlow < 0 {
		headerlow = 0
	}
	if headerhigh > Vmax {
		headerhigh = Vmax
	}
	rowlow := Locx - MaxScanRng
	rowhigh := Locx + MaxScanRng
	if rowlow < 0 {
		rowlow = 0
	}
	if rowhigh > Hmax {
		rowhigh = Hmax
	}

	dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, rowlow, headerlow, rowlow, headerhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)

	return false
}

func showScripts(username string, conn *Conn, usersdb *sql.DB) {
	var InitialCommand string
	usersdb.QueryRow("select InitialCommand from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Initial Command: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Zero from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 0: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select One from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 1: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Two from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 2: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Three from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 3: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Four from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 4: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Five from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 5: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Six from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 6: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Seven from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 7: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Eight from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 8: ")
	sendln(conn, InitialCommand)
	usersdb.QueryRow("select Nine from users WHERE name = ?", username).Scan(&InitialCommand)
	send(conn, "Script 9: ")
	sendln(conn, InitialCommand)

}

//
// Set command with parms
// SEt <SIde | PAssword | PRompt | Output | Iofmt | Initial> <PASSWORD | PROMPT | Medium | Long | Short | ABsolute | RElative | BOth | Empire | Coalition | Command line>
// set prompt x
// set output l/m/s
// set iof ab/re/bo
//
func processSet(comnd string, prevcommand string, invalidcommand bool, conn *Conn, username string, ininit bool, objectsdb *sql.DB, usersdb *sql.DB) bool {
	var InitialCommand string
	var su int
	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	parameter1 := strings.Trim(prm1[1], ctlsp)

	// fmt.Println("aaaa set parameter1:", parameter1," len(prm1):",len(prm1), " prm1:", prm1)

	// If initial is not provided, print the existing command
	if len(prm1) == 2 {
		if testCommandMatch("initial", parameter1, lenOfString2) {
			usersdb.QueryRow("select InitialCommand from users WHERE name = ?", username).Scan(&InitialCommand)
			send(conn, "Current initial command line: ")
			sendln(conn, InitialCommand)
			return false
		}

		// If scripts command, just list them all
		if testCommandMatch("scripts", parameter1, lenOfString2) {
			showScripts(username, conn, usersdb)
			return false
		}
	}

	//
	// Got 2 parms, process Intial - store initial commands to be used in the game, default is none
	//
	if testCommandMatch("initial", parameter1, lenOfString2) {
		// Update users db to have a initial command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command
		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")

		usersdb.Exec("UPDATE users set InitialCommand = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	// fmt.Println("going into 0 parameter1:", parameter1, " prevcommand:", prevcommand, comnd, " testcommandmatch:", testCommandMatch("0", parameter1, lenOfString1))
	//
	// Got 2 parms, process 0 - store 0 command script to be used in the game, default is none
	//
	if testCommandMatch("0", parameter1, lenOfString1) {
		// Update users db to have a 0 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		// fmt.Println("in zero got prevcommand", prevcommand, "comnd:",comnd, "updated value: ", strings.Join(parameter1[2:], " "))
		usersdb.Exec("UPDATE users set Zero = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 1 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("1", parameter1, lenOfString1) {
		// Update users db to have a 2 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set One = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 2 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("2", parameter1, lenOfString1) {
		// Update users db to have a 2 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Two = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 3 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("3", parameter1, lenOfString1) {
		// Update users db to have a 3 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Three = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 4 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("4", parameter1, lenOfString1) {
		// Update users db to have a 4 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Four = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 5 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("5", parameter1, lenOfString1) {
		// Update users db to have a 5 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Five = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 6 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("6", parameter1, lenOfString1) {
		// Update users db to have a 6 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Six = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 7 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("7", parameter1, lenOfString1) {
		// Update users db to have a 7 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Seven = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 8 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("8", parameter1, lenOfString1) {
		// Update users db to have a 8 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Eight = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	//
	// Got 2 parms, process 9 - store 1 command script to be used in the game, default is none
	//
	if testCommandMatch("9", parameter1, lenOfString1) {
		// Update users db to have a 9 command for the user
		// don't use just the parameter, use the entire command line (prevcommand)
		// strip out the set ininitial command

		parameter1 := strings.Split(strings.TrimSpace(strings.Trim(prevcommand, ctlonly)), " ")
		usersdb.Exec("UPDATE users set Nine = ? where name = ?", strings.Join(parameter1[2:], " "), username)
		return true
	}

	// not right number of parms?
	if len(prm1) != 3 {
		sendln(conn, "Invalid parameters")
		return false
	}

	// Got 2 parms, process Side. When changing sides you loose your dock if your docked.
	if testCommandMatch("side", parameter1, lenOfString2) {
		//
		// Must be a super user to use this command (or in initial command)
		//
		usersdb.QueryRow("select SuperUser from users where name=?", username).Scan(&su)
		if su == 0 && ininit != true {
			sendln(conn, "Not a superuser, use set initial side command instead")
			return false
		}

		parameter2 := strings.Trim(prm1[2], ctlsp)
		//	fmt.Println("command: ",parameter1, " len:", lenOfString2, " parameter2:", parameter2, " user:",username)
		if testCommandMatch(NameCoalition, parameter2, lenOfString1) {
			_, err := objectsdb.Exec("UPDATE objects set Side = ?, DockFlag = ? WHERE Nme = ?", SideCoalition, "", username)
			if err != nil {
				log.Fatal(err)
			}
			return false
		} else {
			if testCommandMatch(NameEmpire, parameter2, lenOfString1) {
				_, err := objectsdb.Exec("UPDATE objects set Side = ?, DockFlag = ? WHERE Nme = ?", SideEmpire, "", username)
				if err != nil {
					log.Fatal(err)
				}
				return false
			}
		}
	}

	// Got 2 parms, process IOFMT
	if testCommandMatch("iofmt", parameter1, lenOfString1) {
		parameter2 := strings.Trim(prm1[2], ctlsp)
		if testCommandMatch(NameRelative, parameter2, lenOfString1) {

			_, err := objectsdb.Exec("UPDATE objects set IOFmt = ? WHERE Nme = ?", IOFmtRel, username)
			if err != nil {
				log.Fatal(err)
			}
			return false
		} else {
			if testCommandMatch("absolute", parameter2, lenOfString1) {

				_, err := objectsdb.Exec("UPDATE objects set IOFmt = ? WHERE Nme = ?", IOFmtAbs, username)
				if err != nil {
					log.Fatal(err)
				}
				return false
			} else {
				if testCommandMatch("both", parameter2, lenOfString1) {

					both := IOFmtAbs | IOFmtRel

					_, err := objectsdb.Exec("UPDATE objects set IOFmt = ? WHERE Nme = ?", both, username)
					if err != nil {
						log.Fatal(err)
					}
					return false

				} else {
					sendln(conn, "Invalid parameters")
					return false
				}
			}
		}
	}

	// Got 2 parms, process prompt
	if testCommandMatch("prompt", parameter1, lenOfString2) {
		parameter2 := strings.Trim(prm1[2], ctlsp)
		// Update the user's prompt
		objUpdatePrompt(username, parameter2, objectsdb)
		return false
	}

	// Got 2 parms, process password
	if testCommandMatch("password", parameter1, lenOfString2) {
		parameter2 := strings.Trim(prm1[2], ctlsp)
		// Update the user's password
		// mu.Lock()
		objUpdatePass(username, parameter2, usersdb)
		// mu.Unlock()
		return false
	}

	// Got 2 parms, process output
	if testCommandMatch("output", parameter1, lenOfString1) {
		parameter2 := strings.Trim(prm1[2], ctlsp)
		if testCommandMatch(NameLong, parameter2, lenOfString1) {
			// mu.Lock()
			_, err := objectsdb.Exec("UPDATE objects set OutputLen = ? WHERE Nme = ?", OutLenLong, username)
			// mu.Unlock()
			if err != nil {
				log.Fatal(err)
			}
			return false
		} else {
			if testCommandMatch(NameMedium, parameter2, lenOfString1) {
				// mu.Lock()
				_, err := objectsdb.Exec("UPDATE objects set OutputLen = ? WHERE Nme = ?", OutLenMed, username)
				// mu.Unlock()
				if err != nil {
					log.Fatal(err)
				}
				return false
			} else {
				if testCommandMatch(NameShort, parameter2, lenOfString1) {
					// mu.Lock()
					_, err := objectsdb.Exec("UPDATE objects set OutputLen = ? WHERE Nme = ?", OutLenSh, username)
					// mu.Unlock()
					if err != nil {
						log.Fatal(err)
					}
					return false

				} else {
					sendln(conn, "Invalid parameters")
					return false
				}
			}
		}
	}

	// parm entered is wrong
	sendln(conn, invparm)
	return false
}

//
// Do the default Set command - should be an error without parms
//
func dodefSet(invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB, usersdb *sql.DB) bool {
	var myIOFmt int
	var myPrompt string
	var myOutputLen int

	objectsdb.QueryRow("select IOFmt, prompt, OutputLen from objects where Nme=?", username).Scan(&myIOFmt, &myPrompt, &myOutputLen)
	sendln(conn, "Output settings:")
	// Prompt
	send(conn, "Prompt:\t\t\t")
	send(conn, myPrompt)
	sendln(conn, ">")
	// Output len
	send(conn, "Output length: \t\t")
	if myOutputLen == OutLenLong {
		sendln(conn, NameLong)
	}
	if myOutputLen == OutLenMed {
		sendln(conn, NameMedium)
	}
	if myOutputLen == OutLenSh {
		sendln(conn, NameShort)
	}
	// IOFmt
	send(conn, "Input/Output format: \t")
	if myIOFmt == IOFmtRel {
		sendln(conn, NameRelative)
	}
	if myIOFmt == IOFmtAbs {
		sendln(conn, NameAbsolute)
	}
	if myIOFmt == IOFmtBoth {
		sendln(conn, NameBoth)
	}
	//also show the scripts
	showScripts(username, conn, usersdb)
	return false
}

//
// Shields command
// SHields [{UP|DOwn|<TRansfer number>]
//
func processShields(comnd string, username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var myShld int
	var myShipEnergy int
	var newShipEnergy int
	var newShld int

	// get my info
	objectsdb.QueryRow("select Shld, ShipEnergy from objects where Nme=?", username).Scan(&myShld, &myShipEnergy)

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	parameter1 := strings.Trim(prm1[1], ctlsp)
	if len(prm1) == 2 {
		switch true {
		case testCommandMatch(NameUp, parameter1, lenOfString1):
			if objUpdateShldUp(On, username, objectsdb, conn) == true {
				objUpdateShipEnergy(myShipEnergy-100, username, objectsdb)
				// if tractor beam is on, turn it off (properly)
				endTractor(username, conn, objectsdb)
				sendln(conn, "Shields raised, Captain.")
			}

		case testCommandMatch(NameDown, parameter1, lenOfString1):
			if objUpdateShldUp(Off, username, objectsdb, conn) == true {
				objUpdateShipEnergy(myShipEnergy-100, username, objectsdb)
				sendln(conn, "Shields lowered, Captain.")
			}

		default:
			sendln(conn, invparm)
		}
	} else {
		if len(prm1) == 3 {
			parameter2 := strings.Trim(prm1[2], ctlsp)
			if testCommandMatch(NameTransfer, parameter1, lenOfString1) {
				// Do the logic to transfer to get the max shlds = 2500
				// Value must be numeric adjust to max and min amounts
				amt2transfer, err := strconv.Atoi(parameter2)
				if err == nil {
					//
					// First lower transfer to ensure it doesn't exceed limits
					//
					//	# if transfer too much to shield, reduce transfer amt to max it can be

					if myShld+amt2transfer > InitShield {
						amt2transfer = InitShield - myShld
						//// fmt.Println("1) amt2transfer:", amt2transfer)
					}

					// # if transfer too much from shield, reduce transfer amt to max it can be
					if myShld+amt2transfer < 0 {
						amt2transfer = -myShld
						//// fmt.Println("2) amt2transfer:", amt2transfer)
					}

					// # if transfer too much to ship, reduce transfer amt to max it can be
					if myShipEnergy-amt2transfer > InitEnergy {
						amt2transfer = -(InitEnergy - myShipEnergy)
						//// fmt.Println("3) amt2transfer:", amt2transfer)
					}

					// # if transfer too much from ship, reduce transfer amt to max it can be
					if myShipEnergy-amt2transfer < 0 {
						amt2transfer = myShipEnergy
						//// fmt.Println("4) amt2transfer:", amt2transfer)
					}

					// # update ship and shield amounts
					newShipEnergy = myShipEnergy - amt2transfer
					newShld = amt2transfer + myShld
					//// fmt.Println("a) newShipEnergy:", newShipEnergy)
					//// fmt.Println("b) newShld:", newShld)
					//// fmt.Println("c) amt2transfer:", amt2transfer)

					objUpdateShld(newShld, username, objectsdb)
					objUpdateShipEnergy(newShipEnergy, username, objectsdb)
					//
					// End
					//
				} else {
					sendln(conn, invparm)
					return false
				}

			} else {
				sendln(conn, invparm)
				return false
			}
		}
	}
	return false
}

//
// Do the default Shields command
//
func dodefShields(username string, conn *Conn, objectsdb *sql.DB) bool {
	var myShld int
	var myShldUp int
	// get my info
	objectsdb.QueryRow("select Shld, ShldUp from objects where Nme=?", username).Scan(&myShld, &myShldUp)
	doStatusShields(conn, myShld, myShldUp)
	return false
}

//
// Srscan command
//
func processSrscan(invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB) bool {
	var Nme string
	var Locx int
	var Locy int
	var IOFmt int
	var SeenbyEnemy int
	var Side int

	//	doScan(conn, MaxWarpFactor+1, MaxWarpFactor+1, username, objectsdb)

	// Get users' info for doing scan
	_ = objectsdb.QueryRow("select Nme, Locx, Locy, IOFmt, SeenbyEnemy, Side from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &IOFmt, &SeenbyEnemy, &Side)

	headerlow := Locy - 7
	headerhigh := Locy + 7
	if headerlow < 0 {
		headerlow = 0
	}
	if headerhigh > Vmax {
		headerhigh = Vmax
	}
	rowlow := Locx - 7
	rowhigh := Locx + 7
	if rowlow < 0 {
		rowlow = 0
	}
	if rowhigh > Hmax {
		rowhigh = Hmax
	}

	dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, rowlow, headerlow, rowlow, headerhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)

	return false

}

//
// Do the default Srscan command
//
func dodefSrscan(invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB) bool {
	var Nme string
	var Locx int
	var Locy int
	var IOFmt int
	var SeenbyEnemy int
	var Side int

	//	doScan(conn, MaxWarpFactor+1, MaxWarpFactor+1, username, objectsdb)

	// Get users' info for doing scan
	_ = objectsdb.QueryRow("select Nme, Locx, Locy, IOFmt, SeenbyEnemy, Side from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &IOFmt, &SeenbyEnemy, &Side)

	headerlow := Locy - 7
	headerhigh := Locy + 7
	if headerlow < 0 {
		headerlow = 0
	}
	if headerhigh > Vmax {
		headerhigh = Vmax
	}
	rowlow := Locx - 7
	rowhigh := Locx + 7
	if rowlow < 0 {
		rowlow = 0
	}
	if rowhigh > Hmax {
		rowhigh = Hmax
	}

	dotheScan(conn, username, rowhigh, headerlow, rowhigh, headerhigh, rowlow, headerlow, rowlow, headerhigh, false, IOFmt, Locx, Locy, SeenbyEnemy, Side, objectsdb)

	return false
}

//
// Do the default Status header
//
func doStatusHeader(conn *Conn) {
	sendln(conn, "Status Report:")
	return
}

//
// Do the default Status username
//
func doStatusUsername(conn *Conn, username string) {
	send(conn, "Object:\t\t")
	sendln(conn, username)
	return
}

//
// Do the default Status stardate
//
func doStatusSide(conn *Conn, side int) {
	send(conn, "Side:\t\t")
	switch true {
	case side == SideCoalition:
		sendln(conn, NameCoalition)

	case side == SideEmpire:
		sendln(conn, NameEmpire)

	case side == SideNeutral:
		sendln(conn, NameNeutral)

	case side == SideArcheron:
		sendln(conn, NameArcheron)
	}
	return
}

//
// Do the default Status stardate
//
func doStatusStardate(conn *Conn, stardate int) {
	send(conn, "Stardate:\t")
	sendln(conn, strconv.Itoa(stardate))
	return
}

//
// Do the default Status condition
//
func doStatusCondition(conn *Conn, Stat int, myDockFlag string) {
	send(conn, "Condition:\t")
	if Stat == StatG {
		send(conn, StatGreen)
		if myDockFlag != "" {
			send(conn, " + Docked to ")
			sendln(conn, myDockFlag)
		} else {
			sendln(conn, "")
		}
	} else {
		if Stat == StatY {
			send(conn, StatYellow)
			if myDockFlag != "" {
				send(conn, " + Docked to ")
				sendln(conn, myDockFlag)
			} else {
				sendln(conn, "")
			}
		} else {
			send(conn, StatRed)
			if myDockFlag != "" {
				send(conn, " + Docked to ")
				sendln(conn, myDockFlag)
			} else {
				sendln(conn, "")
			}
		}
	}
	return
}

//
// Do the default Status location
//
func doStatusLocation(conn *Conn, Locx int, Locy int, username string, Nme string, objectsdb *sql.DB) {

	var myLocx int
	var myLocy int
	var myioFmt int

	// Get my loc for determining relative output
	objectsdb.QueryRow("select Locx, Locy, IOFmt from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myioFmt)

	send(conn, "Location:\t")
	if myioFmt == IOFmtAbs || myioFmt == IOFmtBoth {
		send(conn, strconv.Itoa(Locx))
		send(conn, ", ")
		send(conn, strconv.Itoa(Locy))
		send(conn, "  ")
	}
	if myioFmt == IOFmtRel || myioFmt == IOFmtBoth {
		send(conn, strconv.Itoa(Locx-myLocx))
		send(conn, " ")
		sendln(conn, strconv.Itoa(Locy-myLocy))
	} else {
		sendln(conn, " ")
	}
	return
}

//
// Do the default Status torpedoes
//
func doStatusTorp(conn *Conn, PhoTor int) {
	send(conn, "Torpedoes:\t")
	sendln(conn, strconv.Itoa(PhoTor))
	return
}

//
// Do the default Status energy
//
func doStatusEnergy(conn *Conn, ShipEnergy int) {
	send(conn, "Energy left:\t")
	sendln(conn, strconv.Itoa(ShipEnergy))
	return
}

//
// Do the default Status damage
//
func doStatusDamage(conn *Conn, ShipDam int) {
	send(conn, "Damage:\t\t")
	sendln(conn, strconv.Itoa(ShipDam))
	return
}

func calcShields(Shld int) float64 {
	pct := float64(Shld) / float64(InitShield) * 100
	return pct
}

//
// Do the default Status shields - always show shield units in absolute value
//
func doStatusShields(conn *Conn, Shld int, ShldUp int) {
	send(conn, "Shields:\t")
	pct := calcShields(Shld)
	//	pct := float64(Shld) / float64(InitShield) * 100
	j := fmt.Sprintf("%.2f", pct)
	if ShldUp == On {
		send(conn, "+")
	} else {
		send(conn, "-")
	}
	send(conn, j)
	send(conn, "%  ")
	send(conn, strconv.Itoa(Abs(Shld)))
	sendln(conn, " units")
	return
}

//
// Do the default Status radio
//
func doStatusRadio(conn *Conn, Radio int) {
	send(conn, "Radio:\t\t")
	if Radio == 0 {
		sendln(conn, "Off")
	} else {
		sendln(conn, "On")
	}
	return
}

//
// Status command   STatus [{Name|Side|RAdio|SHields|ENgines|STardate|TOrp}] | username | RElative | ABsolute] location]
//
func processStatus(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	var Nme string
	var Side int
	var Locx int
	var Locy int
	var PhoTor int
	var ShipEnergy int
	var ShipDam int
	var Shld int
	var ShldUp int
	var Radio int
	var Stat int
	var mySide int
	var SeenbyEnemy int
	var locx int
	var locy int
	var myLocx int
	var myLocy int
	var myIOFmt int
	var StarDate int
	var myDockFlag string

	//username is currently mine!!
	_ = objectsdb.QueryRow("select Nme, Side, Locx, Locy, PhoTor, ShipEnergy, ShipDam, Shld, ShldUp, Radio, Stat, StarDate, IOFmt, DockFlag from objects WHERE Nme = ?", username).Scan(&Nme, &Side, &Locx, &Locy, &PhoTor, &ShipEnergy, &ShipDam, &Shld, &ShldUp, &Radio, &Stat, &StarDate, &myIOFmt, &myDockFlag)

	prm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")

	if len(prm) == 1 {
		sendln(conn, invparm)
		return true
	} else {
		if len(prm) == 2 {
			switch true {
			case testCommandMatch("name", prm[1], lenOfString2):
				doStatusUsername(conn, username)

			case testCommandMatch("side", prm[1], lenOfString2):
				doStatusSide(conn, Side)

			case testCommandMatch("stardate", prm[1], lenOfString2):
				doStatusStardate(conn, StarDate)

			case testCommandMatch("condition", prm[1], lenOfString1):
				doStatusCondition(conn, Stat, myDockFlag)

			case testCommandMatch("location", prm[1], lenOfString1):
				doStatusLocation(conn, Locx, Locy, username, username, objectsdb) //username used for relative output)

			case testCommandMatch(NameTorp, prm[1], lenOfString1):
				doStatusTorp(conn, PhoTor)

			case testCommandMatch("energy", prm[1], lenOfString1):
				doStatusEnergy(conn, ShipEnergy)

			case testCommandMatch("damage", prm[1], lenOfString1):
				doStatusDamage(conn, ShipDam)

			case testCommandMatch("shields", prm[1], lenOfString2):
				doStatusShields(conn, Shld, ShldUp)

			case testCommandMatch("radio", prm[1], lenOfString1):
				doStatusRadio(conn, Radio)

			default: //  Must be username/relative or absolute
				objectsdb.QueryRow("select Side from objects where Nme=?", username).Scan(&mySide)
				err := objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Nme=?", prm[1]).Scan(&Nme, &SeenbyEnemy)
				if err == nil {
					istrue := SeenbyEnemy & mySide
					if istrue > 0 { //must have been seen to see status!  <= Works!
						dodefStatus("", invalidcommand, conn, objectsdb, username, Nme)
					} else {
						sendln(conn, "Out of range")
						return false
					}
				} else {
					sendln(conn, "No such user!")
					return false
				}
			}

		} else {
			if len(prm) == 3 { // Must be a default address (choose rel/ab from user default)
				pival, err := strconv.Atoi(prm[1])
				if err == nil {
					locx = pival
					pival, err = strconv.Atoi(prm[2])
					if err == nil {
						locy = pival
						objectsdb.QueryRow("select IOFmt, Locx, Locy, Side from objects where Nme=?", username).Scan(&myIOFmt, &myLocx, &myLocy, &mySide)
					}
					if myIOFmt == IOFmtRel || myIOFmt == IOFmtBoth {
						if locx > MaxScanRng || locy > MaxScanRng {
							sendln(conn, OOR)
							return false
						}
						err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", myLocx+locx, myLocy+locy).Scan(&Nme, &SeenbyEnemy)
						if err == nil {
							istrue := SeenbyEnemy & mySide
							if istrue > 0 {
								dodefStatusRel("", invalidcommand, conn, objectsdb, myLocx+locx, myLocy+locy, username)
							} else {
								sendln(conn, OOR)
								return false
							}
						} else {
							sendln(conn, invparm)
							return false
						}
					} else { // absolute
						if (Abs(myLocx)-Abs(locx) > MaxScanRng) || (Abs(myLocy)-Abs(locy) > MaxScanRng) {
							sendln(conn, OOR)
							return false
						}
						err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", locx, locy).Scan(&Nme, &SeenbyEnemy)
						if err == nil {
							istrue := SeenbyEnemy & mySide
							if istrue > 0 {
								dodefStatusRel("", invalidcommand, conn, objectsdb, locx, locy, username)
							} else {
								sendln(conn, OOR)
								return false
							}
						}
					}
				}
			} else {
				if len(prm) == 4 {
					if testCommandMatch(NameRelative, prm[1], lenOfString1) {
						pival, err := strconv.Atoi(prm[2])
						if err == nil {
							locx := pival
							pival, err := strconv.Atoi(prm[3])
							if err == nil {
								locy := pival
								if locx > MaxScanRng || locy > MaxScanRng {
									sendln(conn, OOR)
									return false
								}
								objectsdb.QueryRow("select Locx, Locy, Side from objects where Nme=?", username).Scan(&myLocx, &myLocy, &mySide)
								objectsdb.QueryRow("select SeenbyEnemy from objects where Locx=? AND Locy=?", myLocx+locx, myLocy+locy).Scan(&SeenbyEnemy)
								istrue := SeenbyEnemy & mySide
								if istrue > 0 {
									dodefStatusRel("", invalidcommand, conn, objectsdb, myLocx+locx, myLocy+locy, username)
								} else {
									sendln(conn, OOR)
									return false
								}

							} else {
								sendln(conn, invparm)
								return false
							}
						} else {
							sendln(conn, invparm)
							return false
						}
					} else {
						if testCommandMatch("absolute", prm[1], lenOfString1) {
							pival, err := strconv.Atoi(prm[2])
							if err == nil {
								locx := pival
								pival, err := strconv.Atoi(prm[3])
								if err == nil {
									locy := pival
									objectsdb.QueryRow("select Side from objects where Nme=?", username).Scan(&mySide)

									err = objectsdb.QueryRow("select Nme, SeenbyEnemy from objects where Locx=? AND Locy=?", locx, locy).Scan(&Nme, &SeenbyEnemy)
									if err == nil {
										istrue := SeenbyEnemy & mySide
										if istrue > 0 {
											dodefStatusRel("", invalidcommand, conn, objectsdb, locx, locy, username)
										} else {
											sendln(conn, OOR)
											return false
										}
									}
								} else {
									sendln(conn, invparm)
									return false
								}
							} else {
								sendln(conn, invparm)
								return false
							}
						} else { // use default iofmt setting must be numeric
							sendln(conn, invparm)
							return false
						}
					}
				} else {
					sendln(conn, invparm)
					return false
				}
				return false
			}
		}
	}
	return false
}

//
// Do the default Status command for relative queries
//
func dodefStatusRel(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, locx int, locy int, username string) bool {
	var Nme string
	var Side int
	var Locx int
	var Locy int
	var PhoTor int
	var ShipEnergy int
	var ShipDam int
	var Shld int
	var ShldUp int
	var Radio int
	var Stat int
	var StarDate int
	var myDockFlag string
	var oty int
	var obu int

	_ = objectsdb.QueryRow("select Nme, Side, Locx, Locy, PhoTor, ShipEnergy, ShipDam, Shld, ShldUp, Radio, Stat, StarDate, DockFlag, Objtype, Builds from objects WHERE Locx = ? and Locy = ?", locx, locy).Scan(&Nme, &Side, &Locx, &Locy, &PhoTor, &ShipEnergy, &ShipDam, &Shld, &ShldUp, &Radio, &Stat, &StarDate, &myDockFlag, &oty, &obu)

	doStatusHeader(conn)
	doStatusStardate(conn, StarDate)
	doStatusUsername(conn, Nme)
	if oty == TypePlanet {
		send(conn, "Builds: \t")
		sendln(conn, strconv.Itoa(obu))
	}
	doStatusSide(conn, Side)
	doStatusCondition(conn, Stat, myDockFlag)
	doStatusLocation(conn, Locx, Locy, username, username, objectsdb)
	doStatusTorp(conn, PhoTor)
	doStatusEnergy(conn, ShipEnergy)
	doStatusDamage(conn, ShipDam)
	doStatusShields(conn, Shld, ShldUp)
	doStatusRadio(conn, Radio)
	return false
}

//
// Do the default Status command (no parms)
//
func dodefStatus(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string, Nme string) bool {
	var Side int
	var Locx int
	var Locy int
	var PhoTor int
	var ShipEnergy int
	var ShipDam int
	var Shld int
	var ShldUp int
	var Radio int
	var Stat int
	var StarDate int
	var myIOFmt int
	var myDockFlag string
	var oty int
	var obu int

	_ = objectsdb.QueryRow("select Side, Locx, Locy, PhoTor, ShipEnergy, ShipDam, Shld, ShldUp, Radio, Stat, StarDate, IOFmt, DockFlag, Objtype, Builds from objects WHERE Nme = ?", Nme).Scan(&Side, &Locx, &Locy, &PhoTor, &ShipEnergy, &ShipDam, &Shld, &ShldUp, &Radio, &Stat, &StarDate, &myIOFmt, &myDockFlag, &oty, &obu)

	doStatusHeader(conn)
	doStatusStardate(conn, StarDate)
	doStatusUsername(conn, Nme)
	if oty == TypePlanet {
		send(conn, "Builds: \t")
		sendln(conn, strconv.Itoa(obu))
	}
	doStatusSide(conn, Side)
	doStatusCondition(conn, Stat, myDockFlag)
	doStatusLocation(conn, Locx, Locy, username, Nme, objectsdb)
	doStatusTorp(conn, PhoTor)
	doStatusEnergy(conn, ShipEnergy)
	doStatusDamage(conn, ShipDam)
	doStatusShields(conn, Shld, ShldUp)
	doStatusRadio(conn, Radio)
	return false
}

//
// Do the default Summary Game number
//
func doSummaryGameNumber(conn *Conn) {
	send(conn, "Number of games since system restart:")
	sendln(conn, strconv.Itoa(GameNumber))
	return
}

//
// Do the default Summary Archeron
//
func doSummaryAr(conn *Conn, objectsdb *sql.DB) {
	var count int
	objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideArcheron, TypeShip).Scan(&count)
	send(conn, strconv.Itoa(count))
	count = 0
	sendln(conn, " Archeron ships in the game")
	return
}

//
// Do the default Summary Black holes
//
func doSummaryBH(conn *Conn, objectsdb *sql.DB) {
	var count int
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Objtype = ?", TypeBH).Scan(&count)
	send(conn, strconv.Itoa(count))
	count = 0
	sendln(conn, " Black holes in the game")
	return
}

//
// Do the default Summary Stars
//
func doSummaryStars(conn *Conn, objectsdb *sql.DB) {
	var count int
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Objtype = ?", TypeStar).Scan(&count)
	send(conn, strconv.Itoa(count))
	count = 0
	sendln(conn, " Stars in the game")
	return
}

//
// Do the default Summary Ships
//
func doSummaryShips(conn *Conn, objectsdb *sql.DB) {
	var count int
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideCoalition, TypeShip).Scan(&count)
	send(conn, strconv.Itoa(count))
	count = 0
	sendln(conn, " Coalition ships in the game")
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideEmpire, TypeShip).Scan(&count)
	send(conn, strconv.Itoa(count))
	sendln(conn, " Empire ships in the game")
	doSummaryAr(conn, objectsdb)
	return
}

//
// Do the default Summary Bases
//
func doSummaryBases(conn *Conn, objectsdb *sql.DB) {
	var count int
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideCoalition, TypeBase).Scan(&count)
	send(conn, strconv.Itoa(count))
	count = 0
	sendln(conn, " Coalition bases in the game")
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideEmpire, TypeBase).Scan(&count)
	send(conn, strconv.Itoa(count))
	sendln(conn, " Empire bases in the game")
	count = 0
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side =  ? and Objtype = ?", SideArcheron, TypeBase).Scan(&count)
	if count > 0 {
		send(conn, strconv.Itoa(count))
		sendln(conn, " Archeron bases in the game")
	}
	return
}

//
// Do the default Summary Planets
//
func doSummaryPlanets(conn *Conn, objectsdb *sql.DB) {
	var count int
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideCoalition, TypePlanet).Scan(&count)
	send(conn, strconv.Itoa(count))
	count = 0
	sendln(conn, " Coalition planets in the game")
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideEmpire, TypePlanet).Scan(&count)
	send(conn, strconv.Itoa(count))
	sendln(conn, " Empire planets in the game")
	count = 0
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", SideNeutral, TypePlanet).Scan(&count)
	send(conn, strconv.Itoa(count))
	sendln(conn, " Neutral planets in the game")
	count = 0
	_ = objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side =  ? and Objtype = ?", SideArcheron, TypePlanet).Scan(&count)
	if count > 0 {
		send(conn, strconv.Itoa(count))
		sendln(conn, " Archeron planets in the game")
	}
	return
}

//
// Summary command
//
func processSummary(comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB, username string) bool {
	var Nme string
	var Locx int
	var Locy int
	var PhoTor int
	var ShipEnergy int
	var ShipDam int
	var Shld int
	var Radio int

	_ = objectsdb.QueryRow("select Nme, Locx, Locy, PhoTor, ShipEnergy, ShipDam, Shld, Radio from objects WHERE Nme = ?", username).Scan(&Nme, &Locx, &Locy, &PhoTor, &ShipEnergy, &ShipDam, &Shld, &Radio)

	prm := strings.Split(strings.TrimSpace(strings.ToLower(comnd)), " ")

	if len(prm) == 1 {
		sendln(conn, invparm)
		return true
	}

	if testCommandMatch("ships", prm[1], lenOfString2) == true {
		doSummaryShips(conn, objectsdb)
	} else {
		if testCommandMatch("bases", prm[1], lenOfString1) == true {
			doSummaryBases(conn, objectsdb)
		} else {
			if testCommandMatch("planets", prm[1], lenOfString1) == true {
				doSummaryPlanets(conn, objectsdb)
			} else {
				if testCommandMatch("stars", prm[1], lenOfString2) == true {
					doSummaryStars(conn, objectsdb)
				} else {
					if testCommandMatch("black holes", prm[1], lenOfString1) == true {
						doSummaryBH(conn, objectsdb)
					} else {
						if testCommandMatch("Archeron", prm[1], lenOfString1) == true {
							doSummaryAr(conn, objectsdb)
						} else {
							sendln(conn, invparm)
							return false
						}
					}
				}
			}
		}
	}
	doSummaryGameNumber(conn)
	return false
}

//
// Do the default Summary command
//
func dodefSummary(invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	doSummaryGameNumber(conn)
	sendln(conn, " ")
	doSummaryShips(conn, objectsdb)
	sendln(conn, " ")
	doSummaryBases(conn, objectsdb)
	sendln(conn, " ")
	doSummaryPlanets(conn, objectsdb)
	sendln(conn, " ")
	doSummaryStars(conn, objectsdb)
	sendln(conn, " ")
	doSummaryBH(conn, objectsdb)
	return false
}

//
// Targets command   - no parameters allowed
//
func processTargets(username string, comnd string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	sendln(conn, invparm)
	return false
}

//
// Do the default Targets command
//
func dodefTargets(username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	processList("list target", invalidcommand, conn, objectsdb, username)
	return false
}

//
// Responses from robots (enemys)
//
func doRespRobot(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB) {
	one := [8]string{"Death to", "Destruction to", "I will crush", "Prepare to die,", "You have aroused my wrath,", "You will witness my vengence,", "May you be attacked by a slime-devil,", "I will reduce you to quarks,"}
	two := [8]string{"mindless", "worthless", "ignorant", "idiotic", "stupid", "dumb", "brainless", "Bzzzzzzt"}
	three := [8]string{"mutant", "cretin", "toad", "worm", "parasite", "Vzzzz", "maggot", "alien"}

	send(conn, one[rand.Intn(8)])
	send(conn, " ")
	send(conn, two[rand.Intn(8)])
	send(conn, " ")
	sendln(conn, three[rand.Intn(8)])
}

//
// Responses from Planets
//
func doRespRobotp(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB) {
	one := [8]string{"Received", "Acknowledging", "10-4", "Accepting", "Ignoring", "Recognizing", "Disregarding", "Monitoring"}
	two := [8]string{"information,", "your facts,", "the data,", "your message,", "contact,", "report received,", "interaction,", "transmission,"}
	three := [8]string{"supervisor", "boss", "leader", "commander", "driver", "captain ", "big cheese", "top dog"}

	send(conn, one[rand.Intn(8)])
	send(conn, " ")
	send(conn, two[rand.Intn(8)])
	send(conn, " ")
	sendln(conn, three[rand.Intn(8)])
}

//
// Responses from inanimated objects (black holes, stars)
//
func doRespRoboti(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB) {
	one := [8]string{"Medical officer:", "Communications officer:", "Second in command:", "Starfleet:", "Science officer:", "Cleaning crew:", "Engineering:", "Transporter:"}
	two := [8]string{"We are inquiring on", "We question", "What is wrong with", "Time to check", "Please check", "Please review", "You need to examine", "Are you aware of"}
	three := [8]string{"your health", "your motives", "your mental health", "your physical status", "your senility", "your dementia ", "your mental status", "your weakness"}

	send(conn, one[rand.Intn(8)])
	send(conn, " ")
	send(conn, two[rand.Intn(8)])
	send(conn, " ")
	sendln(conn, three[rand.Intn(8)])
}

//
// Return true if radio is on
//
func checkRadio(ausername string, objectsdb *sql.DB) bool {
	var myRadio int

	// get my info
	objectsdb.QueryRow("select Radio from objects where Nme=?", ausername).Scan(&myRadio)
	if myRadio == On {
		return true
	} else {
		return false
	}
}

//
// Do the send (tell) command
// example of command that has 1 parm and 1 freeform field
// NEEDS side and all handling
// If your radio is off you can't send
//
func doSend(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB, usersdb *sql.DB, pointsdb *sql.DB) bool {
	var myActv int
	var hisActv int
	var hisType int
	var su int
	if username != "" { //I am logged in
		objectsdb.QueryRow("select Actv from objects where Nme=?", username).Scan(&myActv)
		// Can only send messages if not active or (active and radio is on)
		if (checkRadio(username, objectsdb) == true && myActv == On) || (myActv == Off) {
			// mu.Lock()
			sndto := strings.Split(strings.TrimSpace(strings.Trim(comnd, ctlonly)), " ")
			if len(sndto) <= 2 {
				invalidcommand = true
				// mu.Unlock()
				runtime.Gosched()
				invalidcommand = processHelp("help tell", invalidcommand, conn)
				return invalidcommand
			} else {
				invalidcommand = false
			}
			msgtosnd := "Message from " + username + ": " + strings.Join(sndto[2:], " ")

			//
			// Lookup my side
			//
			var mySide int
			var enemySide int
			objectsdb.QueryRow("select Side from objects where Nme=?", username).Scan(&mySide)

			//
			// Figure out enemy side
			//
			if mySide == SideCoalition {
				enemySide = SideEmpire
			} else {
				enemySide = SideCoalition
			}
			//
			// FRiendlies/ENemy/Side and ALL handling
			// Note: if radio of receiver is off, skip sending message! (since = side must be active)
			//
			if testCommandMatch(NameCoalition, sndto[1], lenOfString1) {
				var Nme string
				rows, _ := objectsdb.Query("select Nme from objects where Objtype = ? and Side = ?;", TypeShip, SideCoalition)
				for rows.Next() {
					rows.Scan(&Nme)
					if checkRadio(Nme, objectsdb) {
						sendln(Conmap[Nme].Connection.(*Conn), msgtosnd)
					}
				}
			} else {
				if testCommandMatch(NameEmpire, sndto[1], lenOfString1) {
					var Nme string
					rows, _ := objectsdb.Query("select Nme from objects where Objtype = ? and Side = ?;", TypeShip, SideEmpire)
					for rows.Next() {
						rows.Scan(&Nme)
						if checkRadio(Nme, objectsdb) {
							sendln(Conmap[Nme].Connection.(*Conn), msgtosnd)
						}
					}
				} else { // sending to a robot?  You can get response since your sending (radio must be on)
					if testCommandMatch(NameNeutral, sndto[1], lenOfString1) {
						send(conn, "Neutrals: ")
						doRespRobot("", invalidcommand, conn, username, objectsdb)
					} else { // sending to a robot?
						if testCommandMatch(NameArcheron, sndto[1], lenOfString2) {
							send(conn, "Archeron: ")
							doRespRobot("", invalidcommand, conn, username, objectsdb)
						} else {
							if testCommandMatch(NameAll, sndto[1], lenOfString2) {
								var Nme string
								rows, _ := objectsdb.Query("select Nme from objects where Objtype = ? and (Side = ? or Side = ?);", TypeShip, SideCoalition, SideEmpire)
								for rows.Next() {
									rows.Scan(&Nme)
									if checkRadio(Nme, objectsdb) {
										sendln(Conmap[Nme].Connection.(*Conn), msgtosnd)
									}
								}
							} else {
								if testCommandMatch(NameFriendly, sndto[1], lenOfString1) {
									var Nme string
									rows, _ := objectsdb.Query("select Nme from objects where Objtype = ? and Side = ?;", TypeShip, mySide)
									for rows.Next() {
										rows.Scan(&Nme)
										if checkRadio(Nme, objectsdb) {
											sendln(Conmap[Nme].Connection.(*Conn), msgtosnd)
										}
									}
								} else {
									if testCommandMatch(NameEnemy, sndto[1], lenOfString1) {
										var Nme string
										rows, _ := objectsdb.Query("select Nme from objects where Objtype = ? and Side = ?;", TypeShip, enemySide)
										for rows.Next() {
											rows.Scan(&Nme)
											if checkRadio(Nme, objectsdb) {
												sendln(Conmap[Nme].Connection.(*Conn), msgtosnd)
											}
										}
										// also send robot responses. You can get response since your sending (radio must be on)
										send(conn, "Neutrals: ")
										doRespRobot("", invalidcommand, conn, username, objectsdb)
										send(conn, "Archeron: ")
										doRespRobot("", invalidcommand, conn, username, objectsdb)
									} else {
										//
										// Individual ship handling    zzzzzzzzzzzzz do name matching here use commandmatch
										//
										_, exists := Conmap[sndto[1]]
										if exists {
											objectsdb.QueryRow("select Actv from objects where Nme=?", sndto[1]).Scan(&hisActv)
											if (checkRadio(sndto[1], objectsdb) == true && hisActv == On) || (hisActv == Off) {
												sendln(Conmap[sndto[1]].Connection.(*Conn), msgtosnd)
											} else {
												sendln(conn, "That user's radio is off!")
											}
										} else {
											// still could be a object, get side and respond accordinginly (planet or inatimate object).  If admin, they can send godlike commands to objects.
											usersdb.QueryRow("select SuperUser from users where name=?", username).Scan(&su)
											if su == 0 {
												err := objectsdb.QueryRow("select Actv, Objtype from objects where Nme=?", sndto[1]).Scan(&hisActv, &hisType)
												if hisType == TypePlanet {
													if (checkRadio(sndto[1], objectsdb) == true && hisActv == On) || (hisActv == Off) && err == nil {
														doRespRobotp("", invalidcommand, conn, username, objectsdb)
													} else {
														s := "Invalid user:" + sndto[1]
														sendln(conn, s)
													}
												} else {
													if (checkRadio(sndto[1], objectsdb) == true && hisActv == On) || (hisActv == Off) && err == nil {
														doRespRoboti("", invalidcommand, conn, username, objectsdb)
													} else {
														s := "Invalid user:" + sndto[1]
														sendln(conn, s)
													}
												}
											} else {
//												fmt.Println("in doing command, msg:", strings.Join(sndto[2:], " "), " and to:", sndto[1])
												processCommand(strings.Join(sndto[2:], " "), " ", nil, sndto[1], nil, false, usersdb, objectsdb, pointsdb)
											}
										}
									}
								}
							}
						}
					}
				}
			}
			// mu.Unlock()
			runtime.Gosched()
		} else {
			sendln(conn, "Your radio is off, captain")
		}
	} else {
		sendln(conn, "You must be logged in to send messages")
		invalidcommand = false
	}
	return invalidcommand
}

//
// Time command - show system status
//
func processTime(comnd string, invalidcommand bool, conn *Conn) bool {
	sendln(conn, invparm)
	return false
}

//
// Do the default time command
//
func dodefTime(invalidcommand bool, conn *Conn) bool {
	send(conn, "Server up since:\t")
	sendln(conn, startupTime.Format("Monday, 02-Jan-06 15:04:05 MST"))
	send(conn, "Current time:\t\t")
	sendln(conn, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"))
	send(conn, "Time game started:\t")
	sendln(conn, gameTime.Format("Monday, 02-Jan-06 15:04:05 MST"))
	return false
}

//
// Torpedoes command  TOrpedoes [Absolute|Relative|Computed] #burst [<v1 h1> | name]
// tor locx locy
// tor comp name
// tor 1 locx locy
// tor 2 comp name
// to a/r locx locy
// to 3 a/r  locx locy
// NO checking for accuracy - shoot at space, no prob!
//
func processTorpedoes(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB) bool {
	var myLocx int
	var myLocy int
	var myShldUp int
	var myShipEnergy int
	var myIOFmt int
	var myOutputLen int
	var mySide int
	var myPhoTorDam int
	var myPhoTor int
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var parameter4 string
	var locx int
	var locy int
	var err error
	var err1 error
	var rslt bool
	var pival int
	var numtorps int
	var whoNme string
	var nummoves int
	var incx float32
	var incy float32
	var whoLocx int
	var whoLocy int
	var whoObjtype int
	var whoSeenbyEnemy int
	var whoSide int
	var i int
	var j int
	var objFound bool
	var newlocx int
	var newlocy int
	var phasHit int
	var whoShld int
	var whoShldUp int
	var whoWarpEngDam int
	var whoImpEngDam int
	var whoPhoTorDam int
	var whoPhasDam int
	var whoShldDam int
	var whoCmpDam int
	var whoLifeSupDam int
	var whoRadioDam int
	var whoTractorDam int
	var whoShipDam int

	// Default amount to shoot
	numtorps = defTorpAmt
	//
	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")
	// First parm
	parameter1 = strings.Trim(prm1[1], ctlsp)

	//// fmt.Println("&&&&&  prm1:",prm1, " parameter1;", parameter1)

	// Lookup me
	objectsdb.QueryRow("select Locx, Locy, ShldUp, ShipEnergy, IOFmt, Side, PhoTorDam, PhoTor, OutputLen from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myShipEnergy, &myIOFmt, &mySide, &myPhoTorDam, &myPhoTor, &myOutputLen)
	if len(prm1) == 3 { // 2 parm - tor x y or tor comp name  (default is numtorps)
		parameter2 = strings.Trim(prm1[2], ctlsp)

		// Lookup user
		locx, err = strconv.Atoi(parameter1)
		locy, err1 = strconv.Atoi(parameter2)

		if err == nil && err1 == nil { // numeric parm1&2
			if myIOFmt == IOFmtRel {
				locx = myLocx + locx
				locy = myLocy + locy
			}

		} else { // alpha parm 1?
			if testCommandMatch(NameComputed, parameter1, lenOfString1) {
				locx, locy, rslt = docomputed(username, parameter2, objectsdb)
				if rslt {
					parameter1 = strconv.Itoa(locx)
					parameter2 = strconv.Itoa(locy)
					myIOFmt = IOFmtAbs
				}
			}
		}
	} else {
		if len(prm1) == 4 { // 3 parm - tor 1 locx locy OR tor 2 comp name OR to a/r locx locy
			numtorps, err = strconv.Atoi(parameter1)
			if err == nil { // got an numtorps - must be tor 1 locx locy OR tor 2 comp name
				parameter2 = strings.Trim(prm1[2], ctlsp)
				parameter3 = strings.Trim(prm1[3], ctlsp)
				// Lookup user
				locx, err = strconv.Atoi(parameter2)
				locy, err1 = strconv.Atoi(parameter3)

				// if 2nd parm is alpha, then 3rd must be alpha (for ph com side/name)
				//// fmt.Println("************ err:",err, " err1:",err1, "parameter2:",parameter2, " parameter3:",parameter3)
				if err != nil && err1 != nil {
					if testCommandMatch(NameComputed, parameter2, lenOfString1) {
						//// fmt.Println("got parm2:", parameter2, " len:", lenOfString1, " NameComputed:", NameComputed)
						locx, locy, rslt = docomputed(username, parameter3, objectsdb)
						//// fmt.Println("got locx:", locx, " locy:", locy, " rslt:", rslt, " username:",username, " parameter3:", parameter3)
						if rslt {
							parameter1 = strconv.Itoa(locx)
							parameter2 = strconv.Itoa(locy)
							myIOFmt = IOFmtAbs
						}
					}
				} else { //must be numeric for both values
					if myIOFmt == IOFmtRel || myIOFmt == IOFmtBoth {
						locx = myLocx + locx
						locy = myLocy + locy
					}
				}
			} else {
				numtorps = defTorpAmt
				if testCommandMatch(NameRelative, parameter1, lenOfString1) {
					parameter2 = strings.Trim(prm1[2], ctlsp)
					parameter3 = strings.Trim(prm1[3], ctlsp)
					pival, err = strconv.Atoi(parameter2)
					locx = myLocx + pival
					pival, err1 = strconv.Atoi(parameter3)
					locy = myLocy + pival
					if err != nil || err1 != nil {
						sendln(conn, invparm)
						return false
					}
				} else {
					if testCommandMatch(NameAbsolute, parameter1, lenOfString1) {
						parameter2 = strings.Trim(prm1[2], ctlsp)
						parameter3 = strings.Trim(prm1[3], ctlsp)
						pival, err = strconv.Atoi(parameter2)
						locx = pival
						pival, err1 = strconv.Atoi(parameter3)
						locy = pival
						if err != nil || err1 != nil {
							sendln(conn, invparm)
							return false
						}
					}
				}
			}
		} else {
			if len(prm1) == 5 { // 4 parm - to 3 a/r  locx locy
				numtorps, err = strconv.Atoi(parameter1)
				if err == nil { // got an numtorps
					parameter2 = strings.Trim(prm1[2], ctlsp)
					parameter3 = strings.Trim(prm1[3], ctlsp)
					parameter4 = strings.Trim(prm1[4], ctlsp)
					if testCommandMatch(NameRelative, parameter2, lenOfString1) {
						pival, err = strconv.Atoi(parameter3)
						locx = myLocx + pival
						pival, err1 = strconv.Atoi(parameter4)
						locy = myLocy + pival
						if err != nil || err1 != nil {
							sendln(conn, invparm)
							return false
						}
					} else {
						if testCommandMatch(NameAbsolute, parameter2, lenOfString1) {
							pival, err = strconv.Atoi(parameter3)
							locx = pival
							pival, err1 = strconv.Atoi(parameter4)
							locy = pival
							if err != nil || err1 != nil {
								sendln(conn, invparm)
								return false
							}
						}
					}

				}
			} else {
				sendln(conn, invparm)
				return false
			}
		}
	}

	//
	// Are tubes damaged beyond repair?
	//
	if myPhoTorDam > MaxDam {
		sendln(conn, torpDamMsg)
		return false
	}

	//
	// Is the number of torps <= 0
	//
	if numtorps < 1 {
		sendln(conn, whatCaptain)
		return false
	}

	//
	// Is the number of torps <= 3
	//
	if numtorps > 3 {
		sendln(conn, torptoomany)
		return false
	}

	//
	// Do we have the # torps needed?
	//
	if myPhoTor < numtorps {
		sendln(conn, torpWrongNum)
		return false
	}

	// Calculate the incremental movements
	nummoves, incx, incy = doCalcIncMove(locx-myLocx, locy-myLocy)
	//
	// Got a loc, is it in range?
	//
	if nummoves > MaxScanRng {
		sendln(conn, TOOR)
		return false
	}
	if nummoves <= 0 {
		sendln(conn, invparm)
		return false
	}
	//
	// do torpedo processing here  ************need displacement processing
	//
	for i = 1; i <= numtorps; i++ {
		// reduce your torpedos by 1
		objectsdb.Exec("UPDATE objects set PhoTor = PhoTor - 1 WHERE Nme = ?", username)
		// torpFailure % of torps misfire and hit you
		if rand.Intn(100) > torpFailure { //no misfire
			for j = 1; j <= MaxScanRng; j++ {
				objFound = false
				newlocx = myLocx + int(incx*float32(j))
				newlocy = myLocy + int(incy*float32(j))
				err = objectsdb.QueryRow("select Nme, Locx, Locy, Shld, ShldUp, WarpEngDam, ImpEngDam, PhoTorDam, PhasDam, ShldDam, CmpDam, LifeSupDam, RadioDam, TractorDam, ShipDam, Objtype, SeenbyEnemy, Side from objects WHERE Locx = ? and Locy = ?", newlocx, newlocy).Scan(&whoNme, &whoLocx, &whoLocy, &whoShld, &whoShldUp, &whoWarpEngDam, &whoImpEngDam, &whoPhoTorDam, &whoPhasDam, &whoShldDam, &whoCmpDam, &whoLifeSupDam, &whoRadioDam, &whoTractorDam, &whoShipDam, &whoObjtype, &whoSeenbyEnemy, &whoSide)
				if err == nil { // found something to hit
					// fmt.Println("hit something:", whoNme, whoLocx, whoLocy, whoShld, whoShldUp, whoWarpEngDam, whoImpEngDam, whoPhoTorDam, whoPhasDam, whoShldDam, whoCmpDam, whoLifeSupDam, whoRadioDam, whoTractorDam, whoShipDam, whoObjtype, whoSeenbyEnemy, whoSide)
					if whoSide == mySide {
						notify(username, conn, myLocx, myLocx, msgTorpNeu, whoNme, whoLocx, whoLocy, objectsdb, 0, i)
						return false
					} else {
						phasHit = int(defTorpHitAmt * rand.Float32())
						//						notify(username, conn, myLocx, myLocy, msgTorpHit, whoNme, whoLocx, whoLocy, objectsdb, phasHit, i)
						doHit(whoNme, whoSide, whoLocx, whoLocy, whoShld, whoShldUp, whoWarpEngDam, whoImpEngDam, whoPhoTorDam, whoPhasDam, whoShldDam, whoCmpDam, whoLifeSupDam, whoRadioDam, whoTractorDam, whoShipDam, whoObjtype, phasHit, objectsdb, conn, torpHit, username, myLocx, myLocy, mySide, pointsdb, usersdb)
						objFound = true
						break
					}
				}
			}
			if objFound == false {
				notify(username, conn, myLocx, myLocy, msgTorpMiss, whoNme, whoLocx, whoLocy, objectsdb, phasHit, i)
			}
		} else { //misfire - inflict damage on yourself then stop torps processing
			// gotta do damage only to torps here
			phasHit = int(defTorpHitAmt * rand.Float32())
			notify(username, conn, myLocx, myLocy, msgTorpMisfire, username, myLocx, myLocy, objectsdb, phasHit, 0)
			doHit(username, whoSide, whoLocx, whoLocy, whoShld, whoShldUp, whoWarpEngDam, whoImpEngDam, whoPhoTorDam, whoPhasDam, whoShldDam, whoCmpDam, whoLifeSupDam, whoRadioDam, whoTractorDam, whoShipDam, whoObjtype, phasHit, objectsdb, conn, torpHit, username, myLocx, myLocy, mySide, pointsdb, usersdb)
			return false //A torpedo misfire also aborts the remainder of the burst
		}
	}
	doGameDelay(username, numtorps, objectsdb)
	return false

}

//
// Do the default Torpedoes command
//
func dodefTorpedoes(comnd string, invalidcommand bool, conn *Conn) bool {
	invalidcommand = processHelp("help torp", invalidcommand, conn)
	return false
}

//
// Tractor command
// TractorWho contains who we are tractoring
// Tractoron shows we are BEING tractored
//
func processTractor(comnd string, username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var TractorOn int
	var TractorWho string
	var Nme string
	objectsdb.QueryRow("select TractorOn, TractorWho from objects where Nme=?", username).Scan(&TractorOn, &TractorWho)
	send(conn, "Tractor beam is ")
	if TractorWho != "" {
		send(conn, NameOn)
		send(conn, " to ")
		sendln(conn, TractorWho)
	} else {
		sendln(conn, NameOff)
	}

	if TractorOn == 1 {
		err := objectsdb.QueryRow("select Nme from objects where TractorWho=?", username).Scan(&Nme)
		if err == nil {
			send(conn, Nme)
			sendln(conn, " is tractoring our ship, Captain.")
			return false
		} else {
			sendln(conn, "We are not being tractored, Captain.")
			return false
		}
	} else {
		sendln(conn, "We are not being tractored, Captain.")
	}
	return false
}

//
// Function to end a tractor beam between to players.  Ends only tractorWho connection
//
func endTractorWho(username string, conn *Conn, objectsdb *sql.DB) bool {
	var TractorWho string
	var TractorOn int
	var Objtype int

	//see if we already are tractoring tractored & end
	objectsdb.QueryRow("select TractorWho, TractorOn, Objtype from objects WHERE Nme = ?", username).Scan(&TractorWho, &TractorOn, &Objtype)
	objectsdb.QueryRow("select Objtype from objects WHERE Nme = ?", TractorWho).Scan(&Objtype)

	if TractorWho != "" {
		objectsdb.Exec("UPDATE objects set TractorOn = ? WHERE Nme = ?", 0, TractorWho)
		objectsdb.Exec("UPDATE objects set TractorWho = ? WHERE Nme = ?", "", username)
		if Objtype == TypeShip {
			sendln(conn, endTract)
			sendln(Conmap[TractorWho].Connection.(*Conn), endTract)
		}
	}
	return false
}

//
// Function to end a tractor beam between to players.  Ends tractorOn only
//
func endTractorOn(username string, conn *Conn, objectsdb *sql.DB) bool {
	var TractorWho string
	var TractorOn int
	var Nme string

	//see if we already are tractoring tractored & end
	objectsdb.QueryRow("select TractorWho, TractorOn from objects WHERE Nme = ?", username).Scan(&TractorWho, &TractorOn)

	//see if we already are being tractored & end
	if TractorOn == 1 {
		objectsdb.Exec("UPDATE objects set TractorOn = ? WHERE Nme = ?", 0, username)
		err := objectsdb.QueryRow("select Nme from objects WHERE TractorWho = ?", username).Scan(&Nme)
		if err == nil {
			objectsdb.Exec("UPDATE objects set TractorWho = ? WHERE Nme = ?", "", Nme)
			sendln(Conmap[Nme].Connection.(*Conn), endTract)
		}
		sendln(conn, endTract)
	}
	return false
}

//
// Function to end a tractor beam between to players.  Used for people leaving the game
//
func endTractor(username string, conn *Conn, objectsdb *sql.DB) bool {
	endTractorOn(username, conn, objectsdb)
	endTractorWho(username, conn, objectsdb)
	return false
}

//
// Do the default Tractor command.  TractorOn = I'm being tractored.  TractorWho = who I'm tractoring.
// Rules: if I move with TractorOn, must turn it off and tell whoever is tractoring me, that it's off and remove Tractorwho from them
// If I move with a tractorwho, I have to move them too, and if they have a tractorwho, move their person, etc.
// So first dude will have tractoron=off + a name in tractorwho
// Rest will have tractor on.  Last will have no name in tractorwho.
//		Moves are "follow the leader" .
// 		Move: if tractorOn == 1 turn off, (break tractor if being tractored and move)
// 		Displaced: random turn off tractor/being tractored
// 		Torp/phas hit: random turn off tractor/being tractored
// shields up: turn off tractor and being tractored
// Quit/die: turn off tractor and being tractored
//
// Parms:  	Tractor on name
//         	Tractor on locx locy (for objects) 4
//			Tractor on rel/abs locx locy (for objects) 5
//			tractor off
func dodefTractor(comnd string, username string, invalidcommand bool, conn *Conn, objectsdb *sql.DB) bool {
	var TractorOn int
	var TractorWho string
	var parameter1 string
	var parameter2 string
	var parameter3 string
	var parameter4 string

	// get my info
	objectsdb.QueryRow("select TractorOn, TractorWho from objects where Nme=?", username).Scan(&TractorOn, &TractorWho)

	prm := strings.Trim(comnd, ctlsp)
	prm1 := strings.SplitAfter(prm, " ")

	// First parm
	parameter1 = strings.Trim(prm1[1], ctlsp)

	if len(prm1) == 2 { // 1 parm
		switch true {

		// Turn off tractor - set tractoron to off for tractorwho, then make tractorwho = ""
		case testCommandMatch(NameOff, parameter1, lenOfString2):
			if TractorWho == "" {
				sendln(conn, "Tractor is already off, captain")
			} else {
				objectsdb.Exec("UPDATE objects set TractorOn = ? WHERE Nme = ?", 0, TractorWho)
				TractorWho = ""
				objectsdb.Exec("UPDATE objects set TractorWho = ? WHERE Nme = ?", "", username)
				sendln(conn, "Tractor disengaged, captain")
				return false
			}

		default:
			invalidcommand = processHelp("help tract", invalidcommand, conn)
			return false
		}
	} else { // turn it on to a ship
		if len(prm1) == 3 { // 2 parms
			switch true {
			// parm must be "on shipname"
			case testCommandMatch(NameOn, parameter1, lenOfString2):
				parameter2 = strings.Trim(prm1[2], ctlsp)
				var myLocx int
				var myLocy int
				var myShldUp int
				var myTractorWho string
				var whoLocx int
				var whoLocy int
				var whoShldUp int
				var whoTractorOn int
				//see if we already are tractoring, if shields are up
				objectsdb.QueryRow("select Locx, Locy, ShldUp, TractorWho from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myTractorWho)
				objectsdb.QueryRow("select Locx, Locy, ShldUp, TractorOn from objects WHERE Nme = ?", parameter2).Scan(&whoLocx, &whoLocy, &whoShldUp, &whoTractorOn)
				//
				// Object must be adjacent and shields down for both AND we are not tractoring someone, and they are not being tractored!
				//
				if (Abs(myLocx-whoLocx) <= 1) && (Abs(myLocy-whoLocy) <= 1) {
					if (myShldUp == 0) && (whoShldUp == 0) {
						if (myTractorWho == "") && (whoTractorOn == 0) {
							send(conn, "Tractor engaging:")
							sendln(conn, parameter2)
							objectsdb.Exec("UPDATE objects set TractorWho = ? WHERE Nme = ?", parameter2, username)
							objectsdb.Exec("UPDATE objects set TractorOn = ? WHERE Nme = ?", 1, parameter2)
							return false
						} else {
							sendln(conn, "Either we are already tractoring or they are being tractored, Captain")
							return false
						}
					} else {
						sendln(conn, "Shields must be down for both ships, Captain")
						return false
					}
				} else {
					send(conn, "The ")
					send(conn, parameter2)
					sendln(conn, " is not adjacent")
					return false
				}
			}
		} else {
			if len(prm1) == 4 { // 3 parms must be on locx locy
				switch true {
				// parm must be "on locx locy" <=default iofmt
				case testCommandMatch(NameOn, parameter1, lenOfString2):
					var prm2 int
					var prm3 int
					var err error
					parameter2 = strings.Trim(prm1[2], ctlsp)
					parameter3 = strings.Trim(prm1[3], ctlsp)
					if prm2, err = strconv.Atoi(parameter2); err != nil {
						sendln(conn, invparm)
						return true
					} else {
						if prm3, err = strconv.Atoi(parameter3); err != nil {
							sendln(conn, invparm)
							return true
						}
					}

					var myLocx int
					var myLocy int
					var myShldUp int
					var myTractorWho string
					var whoLocx int
					var whoLocy int
					var whoShldUp int
					var whoTractorOn int
					var whoDat string
					var myIOFmt int
					//see if we already are tractoring, if shields are up
					objectsdb.QueryRow("select IOFmt from objects where Nme=?", username).Scan(&myIOFmt)
					// Am i already tractoring?
					objectsdb.QueryRow("select Locx, Locy, ShldUp, TractorWho from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myTractorWho)

					if myIOFmt == IOFmtAbs {
						objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, TractorOn from objects WHERE Locx = ? and Locy = ?", prm2, prm3).Scan(&whoDat, &whoLocx, &whoLocy, &whoShldUp, &whoTractorOn)
					} else { //must be relative
						objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, TractorOn from objects WHERE Locx = ? and Locy = ?", prm2+myLocx, prm3+myLocy).Scan(&whoDat, &whoLocx, &whoLocy, &whoShldUp, &whoTractorOn)
					}
					//
					// Object must be adjacent and shields down for both AND we are not tractoring someone, and they are not being tractored!
					//
					if (Abs(myLocx-whoLocx) <= 1) && (Abs(myLocy-whoLocy) <= 1) {
						if (myShldUp == 0) && (whoShldUp == 0) {
							if (myTractorWho == "") && (whoTractorOn == 0) {
								send(conn, "Tractor engaging:")
								sendln(conn, whoDat)
								objectsdb.Exec("UPDATE objects set TractorWho = ? WHERE Nme = ?", whoDat, username)
								objectsdb.Exec("UPDATE objects set TractorOn = ? WHERE Nme = ?", 1, whoDat)
								return false
							} else {
								sendln(conn, "Either we are already tractoring or they are being tractored, Captain")
								return false
							}
						} else {
							sendln(conn, "Shields must be down for both objects, Captain")
							return false
						}
					} else {
						sendln(conn, "It is not adjacent")
						return false
					}
				}

			} else {
				if len(prm1) == 5 { // 4 parm must be on rel/ab locx locy

					switch true {
					// parm must be "on relative locx locy" <=default iofmt
					case testCommandMatch(NameOn, parameter1, lenOfString2):
						var myLocx int
						var myLocy int
						var myShldUp int
						var myTractorWho string
						var whoLocx int
						var whoLocy int
						var whoShldUp int
						var whoTractorOn int
						var whoDat string
						var myIOFmt int
						var err error
						var prm3 int
						var prm4 int

						parameter2 = strings.Trim(prm1[2], ctlsp) //rel or abs only
						parameter3 = strings.Trim(prm1[3], ctlsp)
						parameter4 = strings.Trim(prm1[4], ctlsp)

						//see if we already are tractoring, if shields are up
						objectsdb.QueryRow("select IOFmt from objects where Nme=?", username).Scan(&myIOFmt)
						// Am i already tractoring?
						objectsdb.QueryRow("select Locx, Locy, ShldUp, TractorWho from objects WHERE Nme = ?", username).Scan(&myLocx, &myLocy, &myShldUp, &myTractorWho)

						if testCommandMatch(NameRelative, parameter2, lenOfString1) {
							if prm3, err = strconv.Atoi(parameter3); err != nil {
								sendln(conn, invparm)
								return true
							} else {
								if prm4, err = strconv.Atoi(parameter4); err != nil {
									sendln(conn, invparm)
									return true
								}
							}

							objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, TractorOn from objects WHERE Locx = ? and Locy = ?", prm3+myLocx, prm4+myLocy).Scan(&whoDat, &whoLocx, &whoLocy, &whoShldUp, &whoTractorOn)
						} else {
							if testCommandMatch(NameAbsolute, parameter2, lenOfString1) {
								if prm3, err = strconv.Atoi(parameter3); err != nil {
									sendln(conn, invparm)
									return true
								} else {
									if prm4, err = strconv.Atoi(parameter4); err != nil {
										sendln(conn, invparm)
										return true
									}
								}

								objectsdb.QueryRow("select Nme, Locx, Locy, ShldUp, TractorOn from objects WHERE Locx = ? and Locy = ?", prm3, prm4).Scan(&whoDat, &whoLocx, &whoLocy, &whoShldUp, &whoTractorOn)

							} else {
								sendln(conn, invparm)
								return false
							}
						}

						//
						// Object must be adjacent and shields down for both AND we are not tractoring someone, and they are not being tractored!
						//
						if (Abs(myLocx-whoLocx) <= 1) && (Abs(myLocy-whoLocy) <= 1) {
							if (myShldUp == 0) && (whoShldUp == 0) {
								if (myTractorWho == "") && (whoTractorOn == 0) {
									send(conn, "Tractor engaging:")
									sendln(conn, whoDat)
									objectsdb.Exec("UPDATE objects set TractorWho = ? WHERE Nme = ?", whoDat, username)
									objectsdb.Exec("UPDATE objects set TractorOn = ? WHERE Nme = ?", 1, whoDat)
									return false
								} else {
									sendln(conn, "Either we are already tractoring or they are being tractored, Captain")
									return false
								}
							} else {
								sendln(conn, "Shields must be down for both objects, Captain")
								return false
							}
						} else {
							sendln(conn, "It is not adjacent")
							return false
						}
					}

				} else {
					sendln(conn, invparm)
				}
			}
		}
	}
	return false
}

//
// Type command
//
func processType(comnd string, invalidcommand bool, conn *Conn, username string, objectsdb *sql.DB, usersdb *sql.DB) bool {
	var mem runtime.MemStats
	prm := strings.Split(strings.TrimSpace(strings.Trim(comnd, ctlonly)), " ")
	if len(prm) != 2 {
		sendln(conn, invparm)
		return false
	} else {
		if testCommandMatch("system", prm[1], lenOfString1) {
			// Put out time info
			var invalidcommand bool
			dodefTime(invalidcommand, conn)
			//  this stuff is for detailed debugging
			send(conn, "Number of CPUs:\t\t\t\t")
			sendln(conn, strconv.Itoa(runtime.NumCPU()))
			runtime.ReadMemStats(&mem)
			send(conn, "Memory Allocated:\t\t\t")
			sendln(conn, strconv.FormatUint(mem.Alloc, 10))
			send(conn, "Total Memory Allocated:\t\t\t")
			sendln(conn, strconv.FormatUint(mem.TotalAlloc, 10))
			send(conn, "Total bytes from system:\t\t")
			sendln(conn, strconv.FormatUint(mem.Sys, 10))
			send(conn, "Total pointer lookups:\t\t\t")
			sendln(conn, strconv.FormatUint(mem.Lookups, 10))
			send(conn, "Total mallocs:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.Mallocs, 10))
			send(conn, "Number of frees:\t\t\t")
			sendln(conn, strconv.FormatUint(mem.Frees, 10))
			send(conn, "Heap Allocated:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.HeapAlloc, 10))
			send(conn, "Heap Sys:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.HeapSys, 10))

			send(conn, "Heap Idle:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.HeapIdle, 10))
			send(conn, "Heap in use:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.HeapInuse, 10))
			send(conn, "Heap released:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.HeapReleased, 10))
			send(conn, "Heap number of objects:\t\t\t")
			sendln(conn, strconv.FormatUint(mem.HeapObjects, 10))
			send(conn, "Bytes used by stack allocator:\t\t")
			sendln(conn, strconv.FormatUint(mem.StackInuse, 10))
			send(conn, "Bytes used by stack system:\t\t")
			sendln(conn, strconv.FormatUint(mem.StackSys, 10))
			send(conn, "Mspan structures:\t\t\t")
			sendln(conn, strconv.FormatUint(mem.MSpanInuse, 10))
			send(conn, "Mspan sys:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.MSpanSys, 10))
			send(conn, "Mcache structures:\t\t\t")
			sendln(conn, strconv.FormatUint(mem.MCacheInuse, 10))
			send(conn, "Mcache sys:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.MCacheSys, 10))
			send(conn, "Profiling bucket hash table:\t\t")
			sendln(conn, strconv.FormatUint(mem.BuckHashSys, 10))
			send(conn, "GC metadata:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.GCSys, 10))
			send(conn, "Other system allocations:\t\t")
			sendln(conn, strconv.FormatUint(mem.OtherSys, 10))

			send(conn, "Next coll will happen when HeapAlloc ≥:\t")
			sendln(conn, strconv.FormatUint(mem.NextGC, 10))
			send(conn, "End time of last coll in nanoseconds:\t")
			sendln(conn, strconv.FormatUint(mem.LastGC, 10))
			send(conn, "Pause total Ns:\t\t\t\t")
			sendln(conn, strconv.FormatUint(mem.PauseTotalNs, 10))
			//			send(conn, "Circular buffer of recent GC pause durations, most recent at [(NumGC+255)%256]:\t\t\t\t")
			//			sendln(conn, strconv.FormatUint(mem.PauseNs, 10))
			//			send(conn, "Circular buffer of recent GC pause end times:\t\t\t\t")
			//			sendln(conn, strconv.FormatUint(uint64(mem.PauseEnd), 10))
			send(conn, "Other system allocations:\t\t")
			sendln(conn, strconv.FormatUint(uint64(mem.NumGC), 10))
			send(conn, "Fraction of CPU time used by GC:\t")
			sendln(conn, strconv.FormatFloat(mem.GCCPUFraction, 'E', -1, 64))
			send(conn, "Enable GC:\t\t\t\t")
			sendln(conn, strconv.FormatBool(mem.EnableGC))
			send(conn, "Debug GC:\t\t\t\t")
			sendln(conn, strconv.FormatBool(mem.DebugGC))
			send(conn, "Go Version:\t\t\t\t")
			sendln(conn, runtime.Version())
			send(conn, "Number of threads:\t\t\t\t")
			p := pprof.Lookup("threadcreate")
			n := p.Count()
			sendln(conn, strconv.Itoa(n))
			//
			return false
		} else {
			if testCommandMatch("InputOutput", prm[1], lenOfString1) {
				invalidcommand = dodefSet(invalidcommand, conn, username, objectsdb, usersdb)
			} else {
				if testCommandMatch("Game", prm[1], lenOfString1) {
					printVer(conn)
				} else {
					sendln(conn, invparm)
				}
			}
		}
	}
	return false
}

//
// Do the default Type command
//
func dodefType(invalidcommand bool, conn *Conn) bool {
	invalidcommand = processHelp("help type", invalidcommand, conn)
	return false
}

//
// Do the users command
//
func doUsers(invalidcommand bool, conn *Conn) bool {
	err := error(nil)
	// mu.Lock()
	invalidcommand = false
	sendln(conn, "Users logged on:")

	for _, value := range Conmap {
		send(conn, value.Username)
		send(conn, "\t")
		send(conn, value.Remoteaddress)
		send(conn, "\n")
		checkErr(err)
	}
	// mu.Unlock()
	runtime.Gosched()
	return false
}

//
// Process commands
//
func processCommand(comnd string, prevcommand string, conn *Conn, username string, err error, ininit bool, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) string {
	invalidcommand := true
	// Set the user to be active
	objUpdateActv(On, username, objectsdb)
	parameter1 := strings.Split(strings.Trim(comnd, ctlsp), " ")
	//fmt.Println("in command:", parameter1)

	switch true {
	case testCommandMatch("ADministrator", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processAdmin(username, comnd, invalidcommand, conn, conn.Conn.RemoteAddr().String(), usersdb, objectsdb, pointsdb)
		} else {
			invalidcommand = dodefAccount(username, comnd, invalidcommand, conn, conn.Conn.RemoteAddr().String(), usersdb)
		}

	case testCommandMatch("bases", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processBases(comnd, invalidcommand, conn, objectsdb, username)
		} else {
			invalidcommand = dodefBases(invalidcommand, conn, objectsdb, username)
		}

	case testCommandMatch("build", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processBuild(comnd, invalidcommand, conn, objectsdb, username)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefBuild(invalidcommand, conn)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("capture", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processCapture(username, comnd, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefCapture(invalidcommand, conn)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("damages", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processDamages(comnd, username, conn, objectsdb)
		} else {
			invalidcommand = dodefDamages(conn, username, objectsdb)
		}

	case testCommandMatch("dock", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processDock(username, comnd, invalidcommand, conn, objectsdb, false)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefDock(username, comnd, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("energy", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processEnergy(username, comnd, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefEnergy(username, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("gripe", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processGripe(comnd, invalidcommand, conn, username, Conmap[username].Remoteaddress)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefGripe(invalidcommand, conn)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("help", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefHelp(invalidcommand, conn)
		}

		//abbreviate help to a ?
	case testCommandMatch("?", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefHelp(invalidcommand, conn)
		}

		//Do the saved command line for 0
	case testCommandMatch("0", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doZero(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 1
	case testCommandMatch("1", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doOne(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 2
	case testCommandMatch("2", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doTwo(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 3
	case testCommandMatch("3", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doThree(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 4
	case testCommandMatch("4", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doFour(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 5
	case testCommandMatch("5", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doFive(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 6
	case testCommandMatch("6", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doSix(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 7
	case testCommandMatch("7", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doSeven(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 8
	case testCommandMatch("8", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doEight(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

		//Do the saved command line for 9
	case testCommandMatch("9", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = doNine(username, invalidcommand, conn, usersdb, objectsdb, pointsdb)
		}

	case testCommandMatch("impulse", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processImpulse(comnd, invalidcommand, conn, username, objectsdb, pointsdb, usersdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefImpulse(invalidcommand, conn)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("list", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processList(comnd, invalidcommand, conn, objectsdb, username)
		} else {
			invalidcommand = dodefList(invalidcommand, conn, objectsdb, username)
		}

	case testCommandMatch("move", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processMove(comnd, invalidcommand, conn, username, objectsdb, false, pointsdb, usersdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefMove(invalidcommand, conn)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("news", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processNews(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefNews(invalidcommand, conn)
		}

	case testCommandMatch(NamePhas, parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processPhasers(comnd, invalidcommand, conn, objectsdb, username, pointsdb, usersdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefPhasers(invalidcommand, conn)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("planets", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processPlanets(comnd, invalidcommand, conn, objectsdb, username)
		} else {
			invalidcommand = dodefPlanets(username, invalidcommand, conn, objectsdb)
		}

	case testCommandMatch("points", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processPoints(username, comnd, invalidcommand, pointsdb, objectsdb, usersdb, conn)
		} else {
			invalidcommand = dodefPoints(username, invalidcommand, 0, conn, pointsdb, objectsdb, usersdb)
		}

	case testCommandMatch("radio", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processRadio(comnd, username, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefRadio(comnd, username, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("repair", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processRepair(username, comnd, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefRepair(username, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("scan", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processScan(comnd, invalidcommand, conn, username, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefScan(invalidcommand, conn, username, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("set", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processSet(comnd, prevcommand, invalidcommand, conn, username, ininit, objectsdb, usersdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefSet(invalidcommand, conn, username, objectsdb, usersdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("shields", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processShields(comnd, username, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefShields(username, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("srscan", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processSrscan(invalidcommand, conn, username, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefSrscan(invalidcommand, conn, username, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("status", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processStatus(comnd, invalidcommand, conn, objectsdb, username)
		} else {
			invalidcommand = dodefStatus(comnd, invalidcommand, conn, objectsdb, username, username)
		}

	case testCommandMatch(NameSummary, parameter1[0], lenOfString2):
		if testCommandMatch(NameSummary, parameter1[0], lenOfString2) {
			if len(parameter1) > 1 {
				invalidcommand = processSummary(comnd, invalidcommand, conn, objectsdb, username)
			} else {
				invalidcommand = dodefSummary(invalidcommand, conn, objectsdb)
			}
		}

	case testCommandMatch("targets", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processTargets(username, comnd, invalidcommand, conn, objectsdb)
		} else {
			invalidcommand = dodefTargets(username, invalidcommand, conn, objectsdb)
		}

	case testCommandMatch("time", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processTime(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefTime(invalidcommand, conn)
		}

	case testCommandMatch(NameTorp, parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processTorpedoes(comnd, invalidcommand, conn, username, objectsdb, pointsdb, usersdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefTorpedoes(comnd, invalidcommand, conn)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("tractor", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefTractor(comnd, username, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		} else {
			//			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processTractor(comnd, username, invalidcommand, conn, objectsdb)
			//			err = atomicCommit(tx)
			//			runtime.Gosched()
		}

	case testCommandMatch("type", parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processType(comnd, invalidcommand, conn, username, objectsdb, usersdb)
		} else {
			invalidcommand = dodefType(invalidcommand, conn)
		}

	case testCommandMatch("users", parameter1[0], lenOfString1):
		invalidcommand = doUsers(invalidcommand, conn)

	case testCommandMatch("tell", parameter1[0], lenOfString2):
		//		tx, _ := atomicBegin(objectsdb)
		invalidcommand = doSend(comnd, invalidcommand, conn, username, objectsdb, usersdb, pointsdb)
		//		err = atomicCommit(tx)
		//		runtime.Gosched()

	default:
		sendln(conn, "Unknown command -- for help type HELP")

	}
	return username
}

//
// Process Pregame commands
// Activate Gripe Help News Points Quit Summary Time Users
//
//
// Activate command
// Default command no parms (guest name)
// takes username and side as optional parms

func processPregameCommand(comnd string, conn *Conn, username string, c net.Conn, err error, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) string {
	invalidcommand := true

	parameter1 := strings.Split(strings.Trim(comnd, ctlsp), " ")

	objUpdateActv(Off, username, objectsdb)

	switch true {
	case testCommandMatch("activate", parameter1[0], lenOfString1):
		tx, _ := atomicBegin(objectsdb)
		invalidcommand, username = processActivate(comnd, invalidcommand, conn, username, c.RemoteAddr().String(), objectsdb, pointsdb, usersdb)
		err = atomicCommit(tx)
		runtime.Gosched()
		break

	case testCommandMatch("gripe", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			tx, _ := atomicBegin(objectsdb)
			invalidcommand = processGripe(comnd, invalidcommand, conn, username, c.RemoteAddr().String())
			err = atomicCommit(tx)
			runtime.Gosched()
		} else {
			tx, _ := atomicBegin(objectsdb)
			invalidcommand = dodefGripe(invalidcommand, conn)
			err = atomicCommit(tx)
			runtime.Gosched()
		}

	case testCommandMatch("help", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processPregameHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefPregameHelp(invalidcommand, conn)
		}

	case testCommandMatch("?", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processPregameHelp(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefPregameHelp(invalidcommand, conn)
		}

	case testCommandMatch("login", parameter1[0], lenOfString4):
		tx, _ := atomicBegin(objectsdb)
		invalidcommand, username = processLogin(comnd, invalidcommand, conn, conn.Conn.RemoteAddr().String(), username, usersdb)
		err = atomicCommit(tx)
		runtime.Gosched()

	case testCommandMatch("logoff", parameter1[0], lenOfString4):
		tx, _ := atomicBegin(objectsdb)
		invalidcommand, username = processLogoff(username, invalidcommand, conn, conn.Conn.RemoteAddr().String())
		err = atomicCommit(tx)
		runtime.Gosched()

	case testCommandMatch("news", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processNews(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefNews(invalidcommand, conn)
		}

	case testCommandMatch("points", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = dodefPoints(username, invalidcommand, 0, conn, pointsdb, objectsdb, usersdb)
		} else {
			invalidcommand = dodefPoints(username, invalidcommand, 0, conn, pointsdb, objectsdb, usersdb)
		}

	case testCommandMatch("recover", parameter1[0], lenOfString1):
		tx, _ := atomicBegin(objectsdb)
		processRecover(comnd, invalidcommand, conn, conn.Conn.RemoteAddr().String(), usersdb)
		err = atomicCommit(tx)
		runtime.Gosched()

	case testCommandMatch("register", parameter1[0], lenOfString1):
		tx, _ := atomicBegin(objectsdb)
		invalidcommand = processRegister(comnd, invalidcommand, conn, conn.Conn.RemoteAddr().String(), usersdb)
		err = atomicCommit(tx)
		runtime.Gosched()

	case testCommandMatch(NameSummary, parameter1[0], lenOfString2):
		if len(parameter1) > 1 {
			invalidcommand = processSummary(comnd, invalidcommand, conn, objectsdb, username)
		} else {
			invalidcommand = dodefSummary(invalidcommand, conn, objectsdb)
		}

	case testCommandMatch("time", parameter1[0], lenOfString1):
		if len(parameter1) > 1 {
			invalidcommand = processTime(comnd, invalidcommand, conn)
		} else {
			invalidcommand = dodefTime(invalidcommand, conn)
		}

	case testCommandMatch("users", parameter1[0], lenOfString1):
		invalidcommand = doUsers(invalidcommand, conn)

	case testCommandMatch("quit", parameter1[0], lenOfString1):
		tx, _ := atomicBegin(objectsdb)
		invalidcommand = dodefDisconnect(username, invalidcommand, conn)
		err = atomicCommit(tx)
		runtime.Gosched()

	default:
		sendln(conn, "Invalid command")

	}
	return username
}

//
// Delete a record - assume mutex is on already, and data clean
//
func dbDelete(username string) {
	//	delete(Locmap, Plrmap[username].Locn)
	delete(Conmap, username)
}

//
// Insert a record into objects db- assume mutex is on already, and data clean
// Add a record into the points db - ignore is it's already there
//
func dbAddobjects(usersdb *sql.DB, pointsdb *sql.DB, objectsdb *sql.DB, Nme string, Stat int, Actv int, Conndatetime string, Locx int, Locy int,
	SeenbyEnemy int, Objtype int, Side int, WarpEngDam int, ImpEngDam int,
	PhoTorDam int, PhoTor int, PhasDam int, ShldDam int, Shld int, ShldUp int,
	CmpDam int, Cmp int, LifeSupDam int, LifeSup int, RadioDam int, Radio int, TractorDam int,
	TractorOn int, TractorWho string, ShipEnergy int, ShipDam int, Builds int, IOFmt int, OutputLen int, StarDate int, Prompt string, DockFlag string, email string) bool {

	//fmt.Println("dbaddobjects Nme:", Nme, " prompt:",Prompt)
	qry := `INSERT INTO "Objects" (Nme, Stat, Actv, Conndatetime, Locx, Locy,
SeenbyEnemy, Objtype, Side, WarpEngDam, ImpEngDam,
PhoTorDam, PhoTor, PhasDam, ShldDam, Shld, ShldUp,
CmpDam, Cmp, LifeSupDam, LifeSup, RadioDam, Radio, TractorDam,
TractorOn, TractorWho, ShipEnergy, ShipDam, Builds, IOFmt, OutputLen, StarDate, Prompt, DockFlag) VALUES
(` +
		`"` + Nme + `", ` + strconv.Itoa(Stat) + "," + strconv.Itoa(Actv) + "," + `"` + Conndatetime + `",` + strconv.Itoa(Locx) + "," + strconv.Itoa(Locy) + "," +
		strconv.Itoa(SeenbyEnemy) + "," + strconv.Itoa(Objtype) + "," + strconv.Itoa(Side) + "," + strconv.Itoa(WarpEngDam) + "," + strconv.Itoa(ImpEngDam) + "," +
		strconv.Itoa(PhoTorDam) + "," + strconv.Itoa(PhoTor) + "," + strconv.Itoa(PhasDam) + "," + strconv.Itoa(ShldDam) + "," + strconv.Itoa(Shld) + "," + strconv.Itoa(ShldUp) + "," +
		strconv.Itoa(CmpDam) + "," + strconv.Itoa(Cmp) + "," + strconv.Itoa(LifeSupDam) + "," + strconv.Itoa(LifeSup) + "," + strconv.Itoa(RadioDam) + "," + strconv.Itoa(Radio) + "," + strconv.Itoa(TractorDam) + "," +
		strconv.Itoa(TractorOn) + "," + `"",` + strconv.Itoa(ShipEnergy) + "," + strconv.Itoa(ShipDam) + "," + strconv.Itoa(Builds) + "," + strconv.Itoa(IOFmt) + "," + strconv.Itoa(OutputLen) +
		"," + strconv.Itoa(StarDate) + "," + `"` + Prompt + `"` + "," + `""` + ");"

	objectsdb.Exec(qry)
	//	fmt.Println("dbaddobjects error:", erre, " query:",qry)
	//
	// Every object has a pointsdb record
	//

	qry1 := `INSERT INTO "points" (Nme, Side, DamToBH, DamToBases, DamToShip,
	DamToStars, DamToPlanets, NumOfShips, NumOfStarDates) VALUES
	(` +
		`"` + Nme + `", ` + strconv.Itoa(Side) + "," + strconv.Itoa(0) + "," + strconv.Itoa(0) + "," + strconv.Itoa(0) + "," +
		strconv.Itoa(0) + "," + strconv.Itoa(0) + "," + strconv.Itoa(0) + "," + strconv.Itoa(0) + `);`
	//	// fmt.Println("qry1:", qry1)

	pointsdb.Exec(qry1)
	//
	// Increment number of ships (cause points record may already exist)
	//
	pointsdb.Exec("UPDATE points SET NumOfShips = NumOfShips + 1  WHERE Nme = ?", Nme)
	usersdb.Exec("UPDATE users SET NumOfShips = NumOfShips + 1 where name = ?", Nme)
	return true
}

//
// Insert a record into users db- assume mutex is on already, and data clean
//
func dbAddusers(username string, email string, conn net.Conn, addr string, time string, pw string, locv string, loch string, usersdb *sql.DB) error {
	tx, err := usersdb.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO users(name, pswd, addr, tme, DamToBH, DamToBases, DamToShip, DamToStars, DamToPlanets, NumOfShips, NumOfStarDates, RecoveryDateSent, mailAddr, SuperUser, CumlativeWins) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, pw, addr, time, 0, 0, 0, 0, 0, 0, 1, 0, time, email, 0, 0)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

//
// Delete a record into objects db- assume mutex is on already, and data clean
//
func dbDelobjects(conn *Conn, objectsdb *sql.DB, Nme string, pointsdb *sql.DB, usersdb *sql.DB) {
	var Locx int
	var Locy int
	var Objtype int
	var Side int
	var qry string
	var tNme int
	var err error
	//
	// Setup random #
	//
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//
	// if ship, tell their side you died
	//
	objectsdb.QueryRow("select Locx, Locy, Side, Objtype from objects WHERE Nme = ?", Nme).Scan(&Locx, &Locy, &Side, &Objtype)

	notify(Nme, conn, Locx, Locy, msgObjDied, Nme, Locx, Locy, objectsdb, 0, 0)

	//
	// Tell em the final score/tell only the guy dieing, not the robots/etc
	//
	if Conmap[Nme].Connection == conn {
		dodefPoints(Nme, false, 0, conn, pointsdb, objectsdb, usersdb)
	}

	//
	// Archerons don't really die, they just move to new loc and have their info reset (note: they are inactive till reset!)
	//
	if Side == SideArcheron {
		ax := r1.Intn(Vmax)
		ay := r1.Intn(Hmax)
		for {
			err = objectsdb.QueryRow("select Nme from objects WHERE Locx = ? and Locy = ?", ax, ay).Scan(&tNme)
			if err != nil {
				_, err = objectsdb.Exec("UPDATE objects set Locx = ?, Locy = ?, StarDate = 0, WarpEngDam = ?, ImpEngDam = ?, PhoTorDam = ?, PhoTor = ?, PhasDam = ?, ShldDam = ?, Shld = ?, ShldUp = ?,CmpDam = ?, Cmp = ?, LifeSupDam = ?, LifeSup = ?, RadioDam = ?, Radio = ?, TractorDam = ?,ShipEnergy = ?, ShipDam = ?, SeenbyEnemy = ? WHERE Nme = ?", ax, ay, 0, 0, 0, 10, 0, 0, InitShield, On, 0, On, 0, 5, 0, On, 0, InitEnergy, 0, SideArcheron, Nme)
				var t string
				var bx int
				var by int
				err = objectsdb.QueryRow("select Nme, Locx, Locy from objects where Nme = ?", Nme).Scan(&t, &bx, &by)
				break
			}
			ax = r1.Intn(Vmax)
			ay = r1.Intn(Hmax)
		}
		// turn off the active flag until to keep from having wierdness (robots turns it back on correctly)
		objUpdateActv(Off, Nme, objectsdb)
		//
		// increment the ship counter
		//
		pointsdb.Exec("UPDATE points SET NumOfShips = NumOfShips + 1  WHERE Nme = ?", Nme)
	} else {
		qry = `DELETE FROM Objects WHERE Nme = "` + Nme + `"`
		objectsdb.Exec(qry)
	}
}

//
//Do pregame
//
func doPregame(conn *Conn, c net.Conn, username string, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) (string, error) {
	myActv := Off
	err := error(nil)
	prevcommand := ""
	objUpdateActv(Off, username, objectsdb)

	// pregame command processor
	for {
		if username == "" {
			send(conn, "pregame>")
		} else {
			send(conn, username+":PG>")
		}
		comnd, err := input(conn)
		//
		// did he close the connection?
		//
		if err != nil {
			conn.Close()
			if username != "" {
				// mu.Lock()
				dbDelete(strings.Trim(username, ctlsp))
				// mu.Unlock()
				runtime.Gosched()
			}
			return username, err
		}
		//
		//  The line should have no extranious spaces, only 1 space allowed
		//
		str := strings.Replace(comnd, "  ", " ", -1)
		for strings.Compare(comnd, str) != 0 {
			comnd = str
			str = strings.Replace(comnd, "  ", " ", -1)
		}

		//
		// Repeat previous command
		//
		if comnd == "\n" {
			comnd = prevcommand
			sendln(conn, prevcommand)

		} else {
			prevcommand = comnd
		}
		//
		//log all input
		//
		comnd = strings.Replace(comnd, "\n", "", -1)
		log.Print(c.RemoteAddr().String() + "***: " + username + ": " + comnd)
		//
		//Parse commands - seperated by /
		//
		comnd1 := strings.Split(strings.Trim(comnd, ctlonly), "/")

		for i := range comnd1 {
			//
			// not allowed to stack activate command
			//
			username = processPregameCommand(comnd1[i], conn, username, c, err, usersdb, objectsdb, pointsdb)
			//
			// Break out if they activate
			//
			objectsdb.QueryRow("select Actv from objects WHERE Nme = ?", username).Scan(&myActv)
			if myActv == On {
				break
			}
		}
		//
		// Break out if they activate
		//
		objectsdb.QueryRow("select Actv from objects WHERE Nme = ?", username).Scan(&myActv)
		if myActv == On {
			break
		}
	}
	return username, err
}
func printVer(conn *Conn) {
	sendln(conn, "Software version: 1.0.2 				Last updated: 07/07/2025")
}

func doPrintHeader(conn *Conn, objectsdb *sql.DB) {

	sendln(conn, " ")
	sendln(conn, "    ____                                                         ")
	sendln(conn, "   / __ \\___  ______      ______ ___________ _________  ____ ___ ")
	sendln(conn, "  / / / / _ \\/ ___/ | /| / / __ `/ ___/ ___// ___/ __ \\/ __ `__ \\")
	sendln(conn, " / /_/ /  __/ /__ | |/ |/ / /_/ / /  (__  )/ /__/ /_/ / / / / / /")
	sendln(conn, "/_____/\\___/\\___/ |__/|__/\\__,_/_/  /____(_)___/\\____/_/ /_/ /_/ ")
	sendln(conn, " ")

	sendln(conn, "  __________________          _-_               _      _-_      _")
	sendln(conn, "  \\________________|)____.---'---`---.____    _(_).---'---`---.(_)_")
	sendln(conn, "                ||    \\----._________.----/    \\----._________.----/")
	sendln(conn, "                ||     / ,'   `---'               `\\   `]-['   /'")
	sendln(conn, "             ___||_,--'  -._                        `\\.' _ `./'")
	sendln(conn, "            /___          ||(-                        ( (_) )")
	sendln(conn, "                `---._____-'                           `._.'\n")

	/*	sendln(conn, "    _________________       .-------.")
		sendln(conn, "   |________________|)   .-'         `-.")
		sendln(conn, "                ||      /               \\")
		sendln(conn, "                ||     /                 \\")
		sendln(conn, "             ___||_,--|       _.--,       |")
		sendln(conn, "            |_________|______/_.-. |      |")
		sendln(conn, "            |______   |      \\_`-' |      |")
		sendln(conn, "                || `--|        `--'       |")
		sendln(conn, "                ||     \\                 /")
		sendln(conn, "    ____________||___   \\               /")
		sendln(conn, "   |________________|)   `-.         .-'")
		sendln(conn, "                            `-------'\n") */

	sendln(conn, "Website: http://decwars.com				email:decwars@gmail.com")
	printVer(conn)
	sendln(conn, "Copyright 2020 Harris S. Newman Consulting, all rights reserved.")
	sendln(conn, " ")
	sendln(conn, "To quick play enter \"Activate\"")
	sendln(conn, " ")
	sendln(conn, "To login enter \"login <username> <password>.\"")
	sendln(conn, "To register enter enter \"register <username> <email address>.\"")
	sendln(conn, "Use the help command if your new.\n")
	doSummaryGameNumber(conn)
	doSummaryShips(conn, objectsdb)
}

//
//
//
func doGame(conn *Conn, c net.Conn, username string, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) (string, error) {
	var shipEnergy int
	var ShipDam int
	var myLocx int
	var myLocy int
	prevcommand := ""
	myActv := On
	myPrompt := "decwars"
	err := error(nil)
	// set initial prompt
	// mu.Lock()
	objectsdb.QueryRow("select Prompt from objects WHERE Nme = ?", username).Scan(&myPrompt)
	// mu.Unlock()
	//

	for {
		// mu.Lock()
		err = objectsdb.QueryRow("select ShipEnergy, ShipDam, Locx, Locy, Actv from objects WHERE Nme = ?", username).Scan(&shipEnergy, &ShipDam, &myLocx, &myLocy, &myActv)
		// mu.Unlock()
		runtime.Gosched()

		//
		// Not active?
		//
		if myActv == Off {
			// fmt.Println("myactive off *******0")
			return username, err
		}
		//
		// If energy = 0 then you die - back to pregame
		//
		if shipEnergy <= 0 {
			// mu.Lock()
			dbDelobjects(conn, objectsdb, username, pointsdb, usersdb)
			// mu.Unlock()
			runtime.Gosched()
			//
			//log death
			//
			send(conn, username)
			send(conn, " ")
			sendln(conn, "RUNS OUT OF ENERGY!")
			log.Print(username + " Ran out of energy.")
			break
		}

		//
		// if you have more than MaxShipDam for the ship damage, you die
		//
		if ShipDam >= MaxShipDam {
			//
			//Notify/do and log death
			//
			notify(username, conn, myLocx, myLocx, msgDestroyed, username, myLocx, myLocx, objectsdb, 0, 0)
			// mu.Lock()
			dbDelobjects(conn, objectsdb, username, pointsdb, usersdb)
			// mu.Unlock()
			runtime.Gosched()
			log.Print(username + " is destroyed")
			break
		}
		//
		//  If damage is 300 or more units, the device is inoperative
		//
		// mu.Lock()
		dodamagecheck(username, conn, objectsdb)
		// mu.Unlock()
		runtime.Gosched()

		// If energy < 1000 then yellow alert
		if shipEnergy < 1000 {
			send(conn, "E")
			send(conn, strconv.Itoa(shipEnergy))
			send(conn, ">")
			// mu.Lock()
			objUpdateStat(StatY, username, objectsdb)
			// mu.Unlock()
			runtime.Gosched()
		} else {
			// mu.Lock()
			objectsdb.QueryRow("select Prompt from objects WHERE Nme = ?", username).Scan(&myPrompt)
			// mu.Unlock()
			runtime.Gosched()
			send(conn, myPrompt)
			send(conn, ">")
		}

		// Get a command line
		comnd, err := input(conn)

		//
		// did he close the connection? - remove from game/pregame
		//
		if err != nil {
			conn.Close()
			// mu.Lock()
			dbDelete(strings.Trim(username, ctlsp))
			dbDelobjects(conn, objectsdb, username, pointsdb, usersdb)
			// mu.Unlock()
			runtime.Gosched()
			return username, err
		}
		//
		//  The line should have no extranious spaces, only 1 space allowed
		//
		str := strings.Replace(comnd, "  ", " ", -1)
		for strings.Compare(comnd, str) != 0 {
			comnd = str
			str = strings.Replace(comnd, "  ", " ", -1)
		}

		//
		// Repeat previous command line
		//
		if comnd == "\n" {
			comnd = prevcommand
			sendln(conn, prevcommand)

		} else {
			prevcommand = comnd
		}
		//
		//log all input get rid of all new lines dammit
		//
		comnd = strings.Replace(comnd, "\n", "", -1)
		log.Print(c.RemoteAddr().String() + ": " + username + ": " + comnd)

		//
		//Parse commands - seperated by /
		//
		comnd1 := strings.Split(strings.Trim(comnd, ctlonly), "/")
		for i := range comnd1 {
			//
			// Update stardate here
			//
			// mu.Lock()
			objUpdateStarDate(username, objectsdb, usersdb, pointsdb)
			// mu.Unlock()
			runtime.Gosched()
			//
			// check for a quit command
			//
			if testCommandMatch("quit", comnd1[i], lenOfString1) {
				// mu.Lock()
				dbDelobjects(conn, objectsdb, username, pointsdb, usersdb)
				// mu.Unlock()
				runtime.Gosched()
				//
				//log quit
				//
				sendln(conn, "Exiting Decwars...")
				log.Print("Connection for:" + username + " Quit.")
				//
				// Process tractor beam if necessary
				//
				endTractor(username, conn, objectsdb)
				//
				// mu.Lock()
				objUpdateActv(Off, username, objectsdb)
				// mu.Unlock()
				runtime.Gosched()
				break
			}

			// Did you get here unactive (typically game ended)
			// mu.Lock()
			err = objectsdb.QueryRow("select Actv from objects where Nme=?", username).Scan(&myActv)
			// mu.Unlock()
			runtime.Gosched()
			//
			// break out of the loop not active
			//
			if err != nil {
				break
			}
			if myActv == Off {
				// fmt.Println("myactive off *******99")
				return username, err
			}
			//
			// Run the command
			username = processCommand(comnd1[i], prevcommand, conn, username, err, false, usersdb, objectsdb, pointsdb)
			// mu.Lock()
			objectsdb.QueryRow("select Actv from objects where Nme=?", username).Scan(&myActv)
			// mu.Unlock()
			runtime.Gosched()
			if myActv == Off {
				// fmt.Println("myactive off *******1")
				return username, err
			}
		}
		// mu.Lock()
		err = objectsdb.QueryRow("select Actv from objects where Nme=?", username).Scan(&myActv)
		// mu.Unlock()
		runtime.Gosched()
		//
		// break out of the loop not active
		//
		if err != nil {
			break
		}
		if myActv == Off {
			// fmt.Println("myactive off *******2")
			return username, err
		}
	}
	return username, err
}

//
//This is where you handle the flow of incoming connections.
//
func handleConnection(c net.Conn, usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB) {
	username := ""
	myActv := Off

	err := error(nil)
	//
	conn, _ := NewConn(c)

	//
	// Negotiate telnet parms
	//	set to cr/lf (verified)
	conn.SetUnixWriteMode(true)
	//
	doPrintHeader(conn, objectsdb)
	//
	for {
		username, err = doPregame(conn, c, username, usersdb, objectsdb, pointsdb)
		//
		if err == nil {
			if username != "" {
				objectsdb.QueryRow("select Actv from objects WHERE Nme = ?", username).Scan(&myActv)
				if myActv == On {
					//
					// log all connection events
					//
					log.Print("User connected: ", username, " connected from: ", c.RemoteAddr())
					//
					username, err = doGame(conn, c, username, usersdb, objectsdb, pointsdb)
					if err != nil {
						break
					}
				}
			}
		} else {
			break
		}
	}
	//
	conn.Close()
	//
	//log all input
	//
	log.Print("Connection for:" + c.RemoteAddr().String() + ":" + username + " disconnected.")
	//
	// Process tractor beam if necessary
	//
	endTractor(username, conn, objectsdb)
	//
	//	Remove the user from the playermap
	//
	// mu.Lock()
	dbDelete(strings.Trim(username, ctlsp))
	// mu.Unlock()
	runtime.Gosched()
	//
}

//
// Likely location of removal of backspace
//
//Get a line of text from the user
func input(t *Conn) (string, error) {
	inpt, err := t.ReadString('\r')
	// extend timeout
	checkErr(t.SetDeadline(time.Now().Add(timeout)))
	//
	inpt = strings.Trim(inpt, "\r\f")
	return inpt, err
}

//Send a row of text to the client (without ending the line)
func send(t *Conn, s string) {
	if t != nil {
		buf := make([]byte, len(s))
		copy(buf, s)
		_, err := t.Write(buf)
		checkErr(err)
		// extend timeout
		checkErr(t.SetDeadline(time.Now().Add(timeout)))
		//
	}
}

//Send a row of text to the client + EOL
func sendln(t *Conn, s string) {
	if t != nil {
		s += "\n"
		send(t, s)
	}
}

//
// Build the object universe
//
func buildObjects(objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB) {

	/* from decwars:
			knplnt==60		;maximum number of planets
			nstar  = int(51 * ran(0)) * 5 + 100
		nhole  = int(41.0 * ran(0) + 10)
	c--	nplnet = int(20.0 + ran(0) * 61.0)
		nplnet = 60			! ALWAYS insert max. # of planets
	*/

	// Universe limits
	//usize := float64(Vmax * Hmax)

	// Planets and stars calculated from the real decwar numbers...
	//	Planetsmax := Abs(int(usize*(rand.NormFloat64()*.1))) / 2 // 10% of vmax * hmax / 2

	//	Starsmax := Abs(int(usize*(rand.NormFloat64()*.1))) / 2   // 10% of vmax * hmax / 2
	//	Starsmax := Abs(int(51 * ((rand.NormFloat64() * .05) + 10)))

	// BHsmax := Abs(int(usize*(rand.NormFloat64()*.1))) / 2 // 10% of vmax * hmax / 2
	//	BHsmax := Abs(int(41 * ((rand.NormFloat64() * .025) + 10)))

	//	BHsmax := 1 //for testing - no black holes

	//
	// Build Planets
	for i := 0; i < int(Planetsmax); i++ {
		success := false
		for {
			success = dbAddobjects(usersdb, pointsdb, objectsdb, "p"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideNeutral, TypePlanet, SideNeutral, 0, 0, 0, InitPhoTor, 0, 0,
				InitShield, Off, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
			if success == true {
				break
			}
		}
	}

	//
	// Build Stars
	for i := 0; i < int(Starsmax); i++ {
		success := false
		for {
			success = dbAddobjects(usersdb, pointsdb, objectsdb, "s"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideNeutral, TypeStar, SideNeutral, 0, 0, 0, InitPhoTor, 0, 0,
				InitShield, Off, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
			if success == true {
				break
			}
		}
	}

	// Build base - coalition
	for i := 0; i < int(Basesmax); i++ {
		success := false
		for {
			success = dbAddobjects(usersdb, pointsdb, objectsdb, "+b"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideCoalition, TypeBase, SideCoalition, 0, 0, 0, InitPhoTor, 0, 0,
				InitShield, On, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
			if success == true {
				break
			}
		}
	}
	// Build base - empire
	for i := 0; i < int(Basesmax); i++ {
		success := false
		for {
			success = dbAddobjects(usersdb, pointsdb, objectsdb, "-b"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideEmpire, TypeBase, SideEmpire, 0, 0, 0, InitPhoTor, 0, 0,
				InitShield, On, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
			if success == true {
				break
			}
		}
	}

	// Build black holes
	for i := 0; i < int(BHsmax); i++ {
		//	for i := 0; i < int(3); i++ { for testing without bh
		success := false
		for {
			success = dbAddobjects(usersdb, pointsdb, objectsdb, "bh"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideNeutral, TypeBH, SideNeutral, 0, 0, 0, InitPhoTor, 0, 0,
				InitShield, Off, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
			if success == true {
				break
			}
		}
	}

	//
	// Wake the Archeron up
	// Build Archeron ship(s)
	//
	for i := 0; i < int(Archeronmax); i++ {
		success := false
		for {
			success = dbAddobjects(usersdb, pointsdb, objectsdb, "a"+strconv.Itoa(i), StatG, On, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), rand.Intn(Vmax), rand.Intn(Vmax), SideArcheron, TypeShip, SideArcheron, 0, 0, 0, InitPhoTor, 0, 0,
				InitShield, On, 0, 0, 0, InitLifeSup, 0, 1, 0, 0, "", InitEnergy, 0, 0, 0, 0, 0, "", "", "")
			if success == true {
				break
			}
		}
	}
	return
}

//
// opens the DBs and creates table schema
//
func openDBs() (usersdb *sql.DB, objectsdb *sql.DB, pointsdb *sql.DB, robotsdb *sql.DB, planetsdb *sql.DB, basesdb *sql.DB) {
	// Open the user db (persistent)
	usersdb, err := sql.Open("sqlite3", "/home/newmanh/go/src/decwars/users.db")
	checkErr(err)

	// Open the objects db (non-persistent)
	// orig	objectsdb, _ = sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared")
	objectsdb, _ = sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared&_busy_timeout=9999")
	//	qry0 := "attach 'file::memdb1:?mode=memory&cache=shared' as objects;"
	//	objectsdb.Exec(qry0)

	// Open the points db (non-persistent)
	pointsdb, _ = sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared&_busy_timeout=9999")
	//	qry1 := "attach 'file::memdb1:?mode=memory&cache=shared' as points;"
	//	pointsdb.Exec(qry1)

	// Open the robots db (non-persistent)
	//	robotsdb, rerr1 := sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared")
	robotsdb, _ = sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared")
	// fmt.Println("opening robots:", rerr1)

	// Open the planets db (non-persistent)
	//	planetsdb, perr1 := sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared")
	planetsdb, _ = sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared")
	// fmt.Println("opening planets:", perr1)

	// Open the bases db (non-persistent)
	//	basesdb, verr1 := sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared")
	basesdb, _ = sql.Open("sqlite3", "file::memdb1?mode=memory&cache=shared")
	// fmt.Println("opening bases:", verr1)

	// persistant for testing !!!!!!!!!!!
	//	objectsdb, err := sql.Open("sqlite3", "./objects.db")

	//
	// users table schema
	//
	sqlStmt := `
	create table users (name text not null primary key, pswd text, addr text, disabled INTEGER, tme datetime, DamToBH INTEGER, DamToBases INTEGER,
	DamToShip INTEGER, DamToStars INTEGER, DamToPlanets INTEGER, NumOfShips INTEGER, NumOfStarDates INTEGER, RecoveryDateSent date, mailAddr text, 
	SuperUser INTEGER, CumlativeWins INTEGER, InitialCommand text, Zero text, One text, Two text, Three text, Four text, Five text, Six text, 
	Seven Text, Eight Text, Nine Text);
	delete from users;
	`
	_, err = usersdb.Exec(sqlStmt)
	//
	// Create a default superuser - username=admin & hsn and all reserved words
	//
	pw := "a"
	/*
		 * not needed here
		 //
		// Random location for user (needs to check for stomping)
		//
		v := strconv.Itoa(rand.Intn(Vmax))
		h := strconv.Itoa(rand.Intn(Hmax))
	*/
	//
	// Put it into db
	//
	tx, err := usersdb.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO users(name, mailAddr, pswd, disabled, tme, DamToBH, DamToBases, DamToShip, DamToStars, DamToPlanets, NumOfShips, NumOfStarDates, RecoveryDateSent, SuperUser, CumlativeWins) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	disab := 0
	_, err = stmt.Exec("admin", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 1, 0)
	_, err = stmt.Exec("hsn", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 1, 0)
	// Reserved words are disabled - must do all combo's
	disab = 1
	_, err = stmt.Exec("en", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("ene", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("enem", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("enemy", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("fr", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("fri", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("frie", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("frien", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("friend", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("co", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("coa", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("coal", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("coali", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("coalit", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("em", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("emp", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("empi", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("empir", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("empire", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("ne", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("neu", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("neut", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("neutr", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("neutra", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("ar", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("arc", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("arch", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("arche", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	_, err = stmt.Exec("archer", "decwars@gmail.com", GetMD5Hash(pw), disab, time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"), 0, 0, 0, 0, 0, 0, 0, time.Now().Format("02-Jan-06"), 0, 0)
	tx.Commit()
	//

	if err != nil {
		log.Printf("users err on decwars admin user create %q: %s\n", err, sqlStmt)
	}

	//
	// objects table schema
	//
	sqlStmt = `CREATE TABLE objects(Nme TEXT, Stat INTEGER, Actv INTEGER, Conndatetime TEXT, Locx INTEGER, Locy INTEGER,
SeenbyEnemy INTEGER, Objtype INTEGER, Side INTEGER, WarpEngDam INTEGER, ImpEngDam INTEGER,
PhoTorDam INTEGER, PhoTor INTEGER, PhasDam INTEGER, ShldDam INTEGER, Shld INTEGER, ShldUp INTEGER,
CmpDam INTEGER, Cmp INTEGER, LifeSupDam INTEGER, LifeSup INTEGER, RadioDam INTEGER, Radio INTEGER, TractorDam INTEGER,
TractorOn INTEGER, TractorWho Text, ShipEnergy INTEGER, ShipDam INTEGER, Builds INTEGER, IOFmt INTEGER, OutputLen INTEGER, 
StarDate INTEGER, Prompt Text, DockFlag Text, PRIMARY KEY (Locx, Locy));
`
	_, err = objectsdb.Exec(sqlStmt)
	if err != nil {
		log.Printf("##### objects err 1 %q: %s\n", err, sqlStmt)
	}

	//
	// points table schema
	//
	sqlStmt = `CREATE TABLE points(Nme TEXT, Side INTEGER, DamToBH INTEGER, DamToBases INTEGER, DamToShip INTEGER, 
	DamToStars INTEGER, DamToPlanets INTEGER, NumOfShips INTEGER, NumOfStarDates INTEGER);
`
	_, err = pointsdb.Exec(sqlStmt)
	if err != nil {
		log.Printf("##### points err 1 %q: %s\n", err, sqlStmt)
	}

	//
	// Create indexes
	//
	sqlStmt = `create index LocxIndex on objects (Locx);`
	_, err = objectsdb.Exec(sqlStmt)
	if err != nil {
		log.Printf("objects err 1 %q: %s\n", err, sqlStmt)
	}
	sqlStmt = `create index LocyIndex on objects (Locy);`
	_, err = objectsdb.Exec(sqlStmt)
	if err != nil {
		log.Printf("objects err 1 %q: %s\n", err, sqlStmt)
	}
	//
	// Create robots table
	//
	//	sqlStmt = "CREATE TABLE robots AS select hitter.Nme as hitterNme, hitter.Locx as hitterLocx, hitter.Locy as hitterLocy, hitter.Side as hitterSide, hitter.Objtype as hitterType, objects.Nme as objNme, objects.Locx as objLocx, objects.Locy as objLocy, objects.Shld as objShld, objects.ShldUp as objShldUp, objects.WarpEngDam as objWarpEngDam, objects.ImpEngDam as objImpEngDam, objects.PhoTorDam as objPhoTorDam, objects.PhasDam as objPhasDam, objects.ShldDam as objShldDam, objects.CmpDam as objCmpDam, objects.LifeSupDam as objLifeSupDam, objects.RadioDam as objRadioDam, objects.TractorDam as objTractorDam, objects.ShipDam as objShipDam, objects.Side as objSide, objects.Objtype as objType from objects inner join objects hitter where (hitterType = " + strconv.Itoa(TypeShip) + " and hitterSide = " + strconv.Itoa(SideArcheron) + " and objLocx between hitterLocx-" + strconv.Itoa(MaxArcheronRng) + " and hitterLocx+" + strconv.Itoa(MaxArcheronRng) + " and objLocy between hitterLocy-" + strconv.Itoa(MaxArcheronRng) + " and hitterLocy+" + strconv.Itoa(MaxArcheronRng) + ")"
	//	_, err = robotsdb.Exec(sqlStmt)
	//	if err != nil {
	//		log.Printf("##### robots err 1 %q: %s\n", err, sqlStmt)
	//	}

	//
	return usersdb, objectsdb, pointsdb, robotsdb, planetsdb, basesdb
}

//
// Robots function
// Archeron support
// Loop through users, if within range, have them "do their thing"
//
func robots(objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB, robotsdb *sql.DB) {
	var objNme string
	var objLocx int
	var objLocy int
	var objSide int
	var objType int
	var hitterNme string
	var hitterLocx int
	var hitterLocy int
	var hitterSide int
	var hitterType int
	var hitterActv int
	var objShld int
	var objShldUp int
	var objWarpEngDam int
	var objImpEngDam int
	var objPhoTorDam int
	var objPhasDam int
	var objShldDam int
	var objCmpDam int
	var objLifeSupDam int
	var objRadioDam int
	var objTractorDam int
	var objShipDam int
	var qry string
	var qry2 string
	var hitterEnergy int
	var hitterShipDam int
	var tx *sql.Tx
	var inloop bool
	var err1 error

	inloop = false
	//
	// Attach to the robotsdb
	//
	//	qry0 := "attach 'file::memdb1:?mode=memory&cache=shared' as robots;"
	//	qerr1, qerr2 := robotsdb.Exec(qry0)
	// // fmt.Println("attach in robots:", qerr1, qerr2, qry0)
	//
	// Create a list of objects to hit (Archerons)
	//
	for {
		//      // fmt.Println("In robots loop", err1)
		if err1 == nil {
			if inloop == false {
				qry = "CREATE TABLE robots AS select hitter.Nme as hitterNme, hitter.Locx as hitterLocx, hitter.Locy as hitterLocy, hitter.Side as hitterSide, hitter.Objtype as hitterType, objects.Nme as objNme, objects.Locx as objLocx, objects.Locy as objLocy, objects.Shld as objShld, objects.ShldUp as objShldUp, objects.WarpEngDam as objWarpEngDam, objects.ImpEngDam as objImpEngDam, objects.PhoTorDam as objPhoTorDam, objects.PhasDam as objPhasDam, objects.ShldDam as objShldDam, objects.CmpDam as objCmpDam, objects.LifeSupDam as objLifeSupDam, objects.RadioDam as objRadioDam, objects.TractorDam as objTractorDam, objects.ShipDam as objShipDam, objects.Side as objSide, objects.Objtype as objType from objects inner join objects hitter where (hitterType = " + strconv.Itoa(TypeShip) + " and hitterSide = " + strconv.Itoa(SideArcheron) + " and objLocx between hitterLocx-" + strconv.Itoa(MaxArcheronRng) + " and hitterLocx+" + strconv.Itoa(MaxArcheronRng) + " and objLocy between hitterLocy-" + strconv.Itoa(MaxArcheronRng) + " and hitterLocy+" + strconv.Itoa(MaxArcheronRng) + ")"
			} else {
				qry = "insert into robots select hitter.Nme as hitterNme, hitter.Locx as hitterLocx, hitter.Locy as hitterLocy, hitter.Side as hitterSide, hitter.Objtype as hitterType, objects.Nme as objNme, objects.Locx as objLocx, objects.Locy as objLocy, objects.Shld as objShld, objects.ShldUp as objShldUp, objects.WarpEngDam as objWarpEngDam, objects.ImpEngDam asobjImpEngDam, objects.PhoTorDam as objPhoTorDam, objects.PhasDam as objPhasDam, objects.ShldDam as objShldDam, objects.CmpDam as objCmpDam, objects.LifeSupDam as objLifeSupDam, objects.RadioDam as objRadioDam, objects.TractorDam as objTractorDam, objects.ShipDam asobjShipDam, objects.Side as objSide, objects.Objtype as objType from objects inner join objects hitter where (hitterType = " + strconv.Itoa(TypeShip) + " and hitterSide = " + strconv.Itoa(SideArcheron) + " and objLocx between hitterLocx-" + strconv.Itoa(MaxArcheronRng) + " and hitterLocx+" + strconv.Itoa(MaxArcheronRng) + " and objLocy between hitterLocy-" + strconv.Itoa(MaxArcheronRng) + " and hitterLocy+" + strconv.Itoa(MaxArcheronRng) + ")"
			}

			tx, _ = atomicBegin(objectsdb)

			//			verr1, verr2 := robotsdb.Exec(qry)
			robotsdb.Exec(qry)

			atomicCommit(tx)
			runtime.Gosched()

			//// fmt.Println("create or insert of robots:",verr1, verr2, qry)
			//
			// Select data from the robots table
			//
			qry2 = "Select hitterNme, hitterLocx, hitterLocy, hitterSide, hitterType, objNme, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objSide, objType from robots;"

			tx, _ = atomicBegin(objectsdb)

			rows4, err4 := robotsdb.Query(qry2)

			if err4 == nil {
				// // fmt.Println("create or insert of robots:",err4, qry2)

				atomicCommit(tx)
				runtime.Gosched()

				for rows4.Next() {
					//
					// Handle locking of sqlite3 db properly - Archerons
					//
					tx, _ = atomicBegin(objectsdb)

					rows4.Scan(&hitterNme, &hitterLocx, &hitterLocy, &hitterSide, &hitterType, &objNme, &objLocx, &objLocy, &objShld, &objShldUp, &objWarpEngDam, &objImpEngDam, &objPhoTorDam, &objPhasDam, &objShldDam, &objCmpDam, &objLifeSupDam, &objRadioDam, &objTractorDam, &objShipDam, &objSide, &objType)
					// Check to see if we need to kill the Archeron
					objectsdb.QueryRow("Select Actv, ShipDam, ShipEnergy from objects where Nme = ?;", hitterNme).Scan(&hitterActv, &hitterShipDam, &hitterEnergy)
					if hitterEnergy > 0 && hitterShipDam < MaxShipDam {
						if hitterActv == On {
							//
							// There is a archeronPercent chance of the archeron doing a torp, phaser, move (each) or capture
							//  PHASERS:
							s1 := rand.NewSource(time.Now().UnixNano())
							r1 := rand.New(s1)
							a := r1.Intn(100)
							b := r1.Intn(100)
							if a <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									processCommand("ph 200 comp en", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									processCommand("ph 200 comp en", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}
							// archerons can move to planets
							if b <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									processCommand("mo com pla", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									processCommand("mo com pla", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}
							//
							// Torpedoes
							//
							s1 = rand.NewSource(time.Now().UnixNano())
							r1 = rand.New(s1)
							a = r1.Intn(100)
							b = r1.Intn(100)
							if a <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									processCommand("torp 1 comp enem", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									processCommand("torp 1 comp enem", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}
							// capture
							if b <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									processCommand("ca comp plan", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									processCommand("ca comp plan", " ", nil, hitterNme, nil, false, usersdb, objectsdb, pointsdb)
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}

							//
							// Move the Archeron here
							//
							s1 = rand.NewSource(time.Now().UnixNano())
							r1 = rand.New(s1)
							a = r1.Intn(100)
							b = r1.Intn(100)
							if a <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									processMove("move comput ship", false, nil, hitterNme, objectsdb, false, pointsdb, usersdb)
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									processMove("move comput ship", false, nil, hitterNme, objectsdb, false, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}
							// build planets
							if b <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									processMove("bui comp", false, nil, hitterNme, objectsdb, false, pointsdb, usersdb)
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									processMove("bui comp", false, nil, hitterNme, objectsdb, false, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}

						}
					} else {
						dbDelobjects(nil, objectsdb, hitterNme, pointsdb, usersdb)
					}

					atomicCommit(tx)
					runtime.Gosched()
					//
					// Finally, go to sleep
					//
					time.Sleep(150 * time.Millisecond)

				}
				robotsdb.Exec("DELETE from robots")
				inloop = true

			}
		}
		//
		//	Update the actv flag to show they can play again
		//
		objectsdb.Exec("UPDATE objects set Actv = ?  WHERE Side = ?", On, strconv.Itoa(SideArcheron))
	}
}

//
// Planets support
// Loop through users, if within range, have them "do their thing"
//
func planets(objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB, planetsdb *sql.DB) {
	var objNme string
	var objLocx int
	var objLocy int
	var objSide int
	var objType int
	var hitterNme string
	var hitterLocx int
	var hitterLocy int
	var hitterSide int
	var hitterType int
	var hitterActv int
	var objShld int
	var objShldUp int
	var objWarpEngDam int
	var objImpEngDam int
	var objPhoTorDam int
	var objPhasDam int
	var objShldDam int
	var objCmpDam int
	var objLifeSupDam int
	var objRadioDam int
	var objTractorDam int
	var objShipDam int
	var qry string
	var qry2 string
	var hitterEnergy int
	var hitterShipDam int
	var tx *sql.Tx
	var inloop bool
	var err1 error

	inloop = false

	//
	// Attach to the planetsdb
	//
	//	qry0 := "attach 'file::memdb1:?mode=memory&cache=shared' as planets;"
	//	planetsdb.Exec(qry0)
	//
	// Create a list of objects to hit (planets)
	//
	for {

		//fmt.Println("In planet loop", err1)

		if err1 == nil {
			if inloop == false {
				qry = "CREATE TABLE planets AS select hitter.Nme as hitterNme, hitter.Locx as hitterLocx, hitter.Locy as hitterLocy, hitter.Side as hitterSide, hitter.Objtype as hitterType, hitter.Actv as hitterActv, objects.Nme as objNme, objects.Locx as objLocx, objects.Locy as objLocy, objects.Shld as objShld, objects.ShldUp as objShldUp, objects.WarpEngDam as objWarpEngDam, objects.ImpEngDam as objImpEngDam, objects.PhoTorDam as objPhoTorDam, objects.PhasDam as objPhasDam, objects.ShldDam as objShldDam, objects.CmpDam as objCmpDam, objects.LifeSupDam as objLifeSupDam, objects.RadioDam as objRadioDam, objects.TractorDam as objTractorDam, objects.ShipDam as objShipDam, objects.Side as objSide, objects.Objtype as objType from objects inner join objects hitter where (hitterType = " + strconv.Itoa(TypePlanet) + " and hitterSide != objSide) and (objLocx between hitterLocx-" + strconv.Itoa(MaxPlanetRng) + " and hitterLocx+" + strconv.Itoa(MaxPlanetRng) + " and objLocy between hitterLocy-" + strconv.Itoa(MaxPlanetRng) + " and hitterLocy+" + strconv.Itoa(MaxPlanetRng) + ")"
			} else {
				qry = "insert into planets select hitter.Nme as hitterNme, hitter.Locx as hitterLocx, hitter.Locy as hitterLocy, hitter.Side as hitterSide, hitter.Objtype as hitterType, hitter.Actv as hitterActv, objects.Nme as objNme, objects.Locx as objLocx, objects.Locy as objLocy, objects.Shld as objShld, objects.ShldUp as objShldUp, objects.WarpEngDam as objWarpEngDam, objects.ImpEngDam as objImpEngDam, objects.PhoTorDam as objPhoTorDam, objects.PhasDam as objPhasDam, objects.ShldDam as objShldDam, objects.CmpDam as objCmpDam, objects.LifeSupDam as objLifeSupDam, objects.RadioDam as objRadioDam, objects.TractorDam as objTractorDam, objects.ShipDam as objShipDam, objects.Side as objSide, objects.Objtype as objType from objects inner join objects hitter where (hitterType = " + strconv.Itoa(TypePlanet) + " and hitterSide != objSide) and (objLocx between hitterLocx-" + strconv.Itoa(MaxPlanetRng) + " and hitterLocx+" + strconv.Itoa(MaxPlanetRng) + " and objLocy between hitterLocy-" + strconv.Itoa(MaxPlanetRng) + " and hitterLocy+" + strconv.Itoa(MaxPlanetRng) + ")"
			}

			tx, _ = atomicBegin(objectsdb)

			//			verr1, verr2 := planetsdb.Exec(qry)
			planetsdb.Exec(qry)
			atomicCommit(tx)
			runtime.Gosched()

			//fmt.Println("create or insert of planets:",qry)

			//
			// Select data from the bases table
			//
			qry2 = "Select hitterNme, hitterLocx, hitterLocy, hitterActv, hitterSide, hitterType, objNme, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objSide, objType from planets;"

			tx, _ = atomicBegin(objectsdb)
			rows4, err4 := planetsdb.Query(qry2)
			atomicCommit(tx)
			runtime.Gosched()

			//                       // fmt.Println("create or insert of planets:",err4, qry2)

			if err4 == nil {
				//fmt.Println("err4:", err4)
				for rows4.Next() {
					//fmt.Println("In rows4.next")
					//
					// Handle locking of sqlite3 db properly - planets
					//
					tx, _ = atomicBegin(objectsdb)

					rows4.Scan(&hitterNme, &hitterLocx, &hitterLocy, &hitterActv, &hitterSide, &hitterType, &objNme, &objLocx, &objLocy, &objShld, &objShldUp, &objWarpEngDam, &objImpEngDam, &objPhoTorDam, &objPhasDam, &objShldDam, &objCmpDam, &objLifeSupDam, &objRadioDam, &objTractorDam, &objShipDam, &objSide, &objType)

					// Check to see if we need to kill the base
					objectsdb.QueryRow("Select ShipDam, ShipEnergy from objects where Nme = ?;", hitterNme).Scan(&hitterShipDam, &hitterEnergy)

					// only have planets hit ships - 3-31-20
					if objType == TypeShip {
						//fmt.Println("in planets, hitterEnergy:",hitterEnergy, " hitterShipDam:", hitterShipDam, " max ship dam:",MaxShipDam)
						if hitterEnergy > 0 && hitterShipDam < MaxShipDam {
							//
							// There is a archeronPercent chance of the base doing a torp, phaser or move (each)
							//  PHASERS:
							s1 := rand.NewSource(time.Now().UnixNano())
							r1 := rand.New(s1)
							a := r1.Intn(100)
							if a <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									doHit(objNme, objSide, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objType, 200, objectsdb, Conmap[objNme].Connection.(*Conn), phasHit, hitterNme, hitterLocx, hitterLocy, hitterSide, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
									// consume energy from hitter
									objectsdb.Exec("UPDATE objects set ShipEnergy = ShipEnergy - 200  WHERE Nme = ?", hitterNme)
								} else {
									doHit(objNme, objSide, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objType, 200, objectsdb, nil, phasHit, hitterNme, hitterLocx, hitterLocy, hitterSide, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
									// consume energy from hitter
									objectsdb.Exec("UPDATE objects set ShipEnergy = ShipEnergy - 200  WHERE Nme = ?", hitterNme)
								}
							}
							//
							// Torpedoes
							//
							s1 = rand.NewSource(time.Now().UnixNano())
							r1 = rand.New(s1)
							a = r1.Intn(100)
							if a <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									processTorpedoes("torp 1 comp enem", false, nil, hitterNme, objectsdb, pointsdb, usersdb)
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									processTorpedoes("torp 1 comp enem", false, nil, hitterNme, objectsdb, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}
						} else {
							dbDelobjects(nil, objectsdb, hitterNme, pointsdb, usersdb)
						}

					}

					atomicCommit(tx)
					runtime.Gosched()
					//
					// Finally, go to sleep
					//
					time.Sleep(150 * time.Millisecond)
				}
			}
		}
		//
		// Delete all recs and start over
		//
		planetsdb.Exec("DELETE from planets")
		inloop = true
		//
		// Slow down planets - have them run only on 3 sec timeframe
		//
		time.Sleep(2500 * time.Millisecond)
	}
}

//
// Bases support
// Loop through users, if within range, have them "do their thing"
//
func bases(objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB, basesdb *sql.DB) {
	var objNme string
	var objLocx int
	var objLocy int
	var objSide int
	var objType int
	var hitterNme string
	var hitterLocx int
	var hitterLocy int
	var hitterSide int
	var hitterType int
	var hitterActv int
	var objShld int
	var objShldUp int
	var objWarpEngDam int
	var objImpEngDam int
	var objPhoTorDam int
	var objPhasDam int
	var objShldDam int
	var objCmpDam int
	var objLifeSupDam int
	var objRadioDam int
	var objTractorDam int
	var objShipDam int
	var qry string
	var qry2 string
	var hitterEnergy int
	var hitterShipDam int
	var tx *sql.Tx
	var inloop bool
	var err1 error
	//	var err2 error

	inloop = false

	//
	// Attach to the basesdb
	//
	//	qry0 := "attach 'file::memdb1:?mode=memory&cache=shared' as bases;"
	//	basesdb.Exec(qry0)
	//
	// Create a list of objects to hit (bases)
	//
	for {
		// fmt.Println("In bases loop")

		if err1 == nil {
			if inloop == false {
				qry = "CREATE TABLE bases AS select hitter.Nme as hitterNme, hitter.Locx as hitterLocx, hitter.Locy as hitterLocy, hitter.Side as hitterSide, hitter.Objtype as hitterType, hitter.Actv as hitterActv, objects.Nme as objNme, objects.Locx as objLocx, objects.Locy as objLocy, objects.Shld as objShld, objects.ShldUp as objShldUp, objects.WarpEngDam as objWarpEngDam, objects.ImpEngDam as objImpEngDam, objects.PhoTorDam as objPhoTorDam, objects.PhasDam as objPhasDam, objects.ShldDam as objShldDam, objects.CmpDam as objCmpDam, objects.LifeSupDam as objLifeSupDam, objects.RadioDam as objRadioDam, objects.TractorDam as objTractorDam, objects.ShipDam as objShipDam, objects.Side as objSide, objects.Objtype as objType from objects inner join objects hitter where (hitterType = " + strconv.Itoa(TypeBase) + " and objects.objType = " + strconv.Itoa(TypeShip) + " and hitterSide != objSide) and (objLocx between hitterLocx-" + strconv.Itoa(MaxBaseRng) + " and hitterLocx+" + strconv.Itoa(MaxBaseRng) + " and objLocy between hitterLocy-" + strconv.Itoa(MaxBaseRng) + " and hitterLocy+" + strconv.Itoa(MaxBaseRng) + ")"
			} else {
				qry = "insert into bases select hitter.Nme as hitterNme, hitter.Locx as hitterLocx, hitter.Locy as hitterLocy, hitter.Side as hitterSide, hitter.Objtype as hitterType, hitter.Actv as hitterActv, objects.Nme as objNme, objects.Locx as objLocx, objects.Locy as objLocy, objects.Shld as objShld, objects.ShldUp as objShldUp, objects.WarpEngDam as objWarpEngDam, objects.ImpEngDam as objImpEngDam, objects.PhoTorDam as objPhoTorDam, objects.PhasDam as objPhasDam, objects.ShldDam as objShldDam, objects.CmpDam as objCmpDam, objects.LifeSupDam as objLifeSupDam, objects.RadioDam as objRadioDam, objects.TractorDam as objTractorDam, objects.ShipDam as objShipDam, objects.Side as objSide, objects.Objtype as objType from objects inner join objects hitter where (hitterType = " + strconv.Itoa(TypeBase) + " and objects.objType = " + strconv.Itoa(TypeShip) + " and hitterSide != objSide) and (objLocx between hitterLocx-" + strconv.Itoa(MaxBaseRng) + " and hitterLocx+" + strconv.Itoa(MaxBaseRng) + " and objLocy between hitterLocy-" + strconv.Itoa(MaxBaseRng) + " and hitterLocy+" + strconv.Itoa(MaxBaseRng) + ")"
			}

			tx, _ = atomicBegin(objectsdb)

			// verr1, verr2 := basesdb.Exec(qry)
			basesdb.Exec(qry)
			atomicCommit(tx)
			runtime.Gosched()

			// fmt.Println("create or insert of bases:",verr1, verr2, qry)
			//
			// Select data from the bases table
			//
			qry2 = "Select hitterNme, hitterLocx, hitterLocy, hitterActv, hitterSide, hitterType, objNme, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objSide, objType from bases;"

			tx, _ = atomicBegin(objectsdb)
			rows4, err4 := basesdb.Query(qry2)
			atomicCommit(tx)
			runtime.Gosched()

			// fmt.Println("select in bases:",err4, qry2)
			if err4 == nil {
				for rows4.Next() {
					// fmt.Println("Rows-next in bases")
					//
					// Handle locking of sqlite3 db properly - bases
					//
					tx, _ = atomicBegin(objectsdb)

					_ = rows4.Scan(&hitterNme, &hitterLocx, &hitterLocy, &hitterActv, &hitterSide, &hitterType, &objNme, &objLocx, &objLocy, &objShld, &objShldUp, &objWarpEngDam, &objImpEngDam, &objPhoTorDam, &objPhasDam, &objShldDam, &objCmpDam, &objLifeSupDam, &objRadioDam, &objTractorDam, &objShipDam, &objSide, &objType)
					// fmt.Println("rows4.Scan err2=",err2, "obj:", objNme, " hitter:",hitterNme)
					// Check to see if we need to kill the base
					_ = objectsdb.QueryRow("Select ShipDam, ShipEnergy from objects where Nme = ?;", hitterNme).Scan(&hitterShipDam, &hitterEnergy)
					// fmt.Println("QueryRow err2=",err2, " hitternme: ",hitterNme, "hitter ship dam: ", hitterShipDam, "hitter energy: ", hitterEnergy, " objtype:", objType, " objectname:", objNme)

					if objType == TypeShip {

						// fmt.Println("hitterenergy:", hitterEnergy, " hitter ship dam:", hitterShipDam)
						if hitterEnergy > 0 && hitterShipDam < MaxBaseDam {
							//
							// There is a archeronPercent chance of the base doing a torp, phaser or move (each)
							//  PHASERS:
							s1 := rand.NewSource(time.Now().UnixNano())
							r1 := rand.New(s1)
							a := r1.Intn(100)
							//						if a <= archeronPercent {

							if true { //testing, should use above
								if Conmap[objNme].Connection != nil {
									// fmt.Println("doing hit 1 hitter:", hitterNme, " object:", objNme)
									doHit(objNme, objSide, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objType, 200, objectsdb, Conmap[objNme].Connection.(*Conn), phasHit, hitterNme, hitterLocx, hitterLocy, hitterSide, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
									// consume energy from hitter
									objectsdb.Exec("UPDATE objects set ShipEnergy = ShipEnergy - 200  WHERE Nme = ?", hitterNme)
								} else {
									// fmt.Println("doing hit 2 hitter:", hitterNme, " object:", objNme)
									doHit(objNme, objSide, objLocx, objLocy, objShld, objShldUp, objWarpEngDam, objImpEngDam, objPhoTorDam, objPhasDam, objShldDam, objCmpDam, objLifeSupDam, objRadioDam, objTractorDam, objShipDam, objType, 200, objectsdb, nil, phasHit, hitterNme, hitterLocx, hitterLocy, hitterSide, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
									// consume energy from hitter
									objectsdb.Exec("UPDATE objects set ShipEnergy = ShipEnergy - 200  WHERE Nme = ?", hitterNme)
								}
							}
							//
							// Torpedoes
							//
							s1 = rand.NewSource(time.Now().UnixNano())
							r1 = rand.New(s1)
							a = r1.Intn(100)
							if a <= archeronPercent {
								if Conmap[objNme].Connection != nil {
									// fmt.Println("doing torps 1 hitter:", hitterNme, " object:", objNme)
									processTorpedoes("torp 1 comp enem", false, nil, hitterNme, objectsdb, pointsdb, usersdb)
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								} else {
									// fmt.Println("doing torps 2 hitter:", hitterNme, " object:", objNme)
									processTorpedoes("torp 1 comp enem", false, nil, hitterNme, objectsdb, pointsdb, usersdb)
									//
									// Update the stardate for the sucker
									//
									objUpdateStarDate(hitterNme, objectsdb, usersdb, pointsdb)
								}
							}
						} else {
							dbDelobjects(nil, objectsdb, hitterNme, pointsdb, usersdb)
						}

					}

					atomicCommit(tx)
					runtime.Gosched()
					//
					// Finally, go to sleep
					//
					time.Sleep(150 * time.Millisecond)

				}
			}
		}
		// fmt.Println("Now deleting bases data")
		//
		// Delete all recs and start over
		//
		basesdb.Exec("DELETE from bases")
		inloop = true
		//
		// Slow down bases - have them run only on 3 sec timeframe
		//
		time.Sleep(2500 * time.Millisecond)
	}
}

//
// Game ender - check for game being over, if so, pop everyone into pregame and reset board
//
func gamender(objectsdb *sql.DB, pointsdb *sql.DB, usersdb *sql.DB) {
	var tx *sql.Tx
	var coBases int
	var emBases int
	var coplnts int
	var emplnts int
	var qry string
	var Nme string

	coBases = 0
	emBases = 0
	for {
		//
		// Step one: Lock the database
		//
		tx, _ = atomicBegin(objectsdb)
		//
		//	Step two: See if there are any ports left for each side, Objtype 8 is base 2 is planets
		//
		objectsdb.QueryRow("Select count(*) from objects where Side = ? and Objtype = ?;", strconv.Itoa(SideCoalition), strconv.Itoa(TypeBase)).Scan(&coBases)
		objectsdb.QueryRow("Select count(*) from objects where Side = ? and Objtype = ?;", strconv.Itoa(SideEmpire), strconv.Itoa(TypeBase)).Scan(&emBases)
		objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", strconv.Itoa(SideCoalition), strconv.Itoa(TypePlanet)).Scan(&coplnts)
		objectsdb.QueryRow("SELECT Count(*) FROM objects WHERE Side = ? and Objtype = ?", strconv.Itoa(SideEmpire), strconv.Itoa(TypePlanet)).Scan(&emplnts)
		//
		// If either side has zero bases and there are no planets, announce game over and reset
		//
		//		fmt.Println("cobases:", coBases, " coplnts:", coplnts, " emBases:", emBases, " emplnts:", emplnts)
		if (coBases == 0 && coplnts == 0) || (emBases == 0 && emplnts == 0) {
			//
			// Reset the game start date/time
			//
			gameTime = time.Now()
			//
			notify("", nil, 0, 0, msgEndGame, "", 0, 0, objectsdb, 0, 0)
			if coBases == 0 {
				notify("", nil, 0, 0, msgEndEmpire, "", 0, 0, objectsdb, 0, 0)
				//
				// Update users to show Empire won one more
				//
				usersdb.Exec("UPDATE users set CumlativeWins = CumlativeWins + 1  WHERE name = ?", "empire")
			} else {
				notify("", nil, 0, 0, msgEndCoalition, "", 0, 0, objectsdb, 0, 0)
				//
				// Update users to show Coalition won one more
				//
				usersdb.Exec("UPDATE users set CumlativeWins = CumlativeWins + 1  WHERE name = ?", "coalit")
			}
			//
			// First loop through all users and kill em
			//
			qry = "select Nme from objects where Actv == 1"
			rows, err := objectsdb.Query(qry)
			if err == nil {
				for rows.Next() {
					rows.Scan(&Nme)
					objUpdateActv(Off, Nme, objectsdb)
				}
			}
			//
			// Now everyone is out of the game, reset the game
			//
			qry = "DELETE FROM Objects"
			objectsdb.Exec(qry)
			//
			// Delete game points db
			//
			qry = "DELETE FROM points"
			pointsdb.Exec(qry)
			//
			// build the universe
			//
			buildObjects(objectsdb, pointsdb, usersdb)
			GameNumber = GameNumber + 1
		}
		//
		// Step three: release the database
		//
		atomicCommit(tx)
		runtime.Gosched()
		//
		// Sleep for a while
		time.Sleep(150 * time.Millisecond)
	}
}

//
//Generic TCP-Server startup
// Now setup to be able to call with parameters:
// Order and default:
// Vmax = 74
// Hmax = 74
// Starsmax = 120
// BHsmax = 120
// Planetsmax = 60
// Archeronmax = 5
// Portnum = 1701

//
func main() {
	var err error
	//
	// First, get the parameters and setup the defaults
	//
	argsWithoutProg := os.Args[1:]

//fmt.Println("len os.Args[1]:", len(os.Args[1:]))
	if len(os.Args[1:]) == 7 {
//fmt.Println("inside num = 7")
		a1 := strings.Trim(argsWithoutProg[0], ctlsp)
		a2 := strings.Trim(argsWithoutProg[1], ctlsp)
		a3 := strings.Trim(argsWithoutProg[2], ctlsp)
		a4 := strings.Trim(argsWithoutProg[3], ctlsp)
		a5 := strings.Trim(argsWithoutProg[4], ctlsp)
		a6 := strings.Trim(argsWithoutProg[5], ctlsp)
		a7 := strings.Trim(argsWithoutProg[6], ctlsp)
//fmt.Println("starting with parms:",a1, a2, a3, a4, a5, a6, a7)
		Vmax, err = strconv.Atoi(a1)
		if err != nil {
			Vmax = 74
		}
		Hmax, err = strconv.Atoi(a2)
		if err != nil {
			Hmax = 74
		}

		Starsmax, err = strconv.Atoi(a3)
		if err != nil {
			Starsmax = 120
		}
		BHsmax, err = strconv.Atoi(a4)
		if err != nil {
			BHsmax = 120
		}
		Archeronmax, err = strconv.Atoi(a5)
		if err != nil {
			Archeronmax = 60
		}
		Planetsmax, err = strconv.Atoi(a6)
		if err != nil {
			Planetsmax = 60
		}

		Portnum, err = strconv.Atoi(a7)
		if err != nil {
		}

	}
	//
	// Setup all syslog output to a file
	//
	logFile, _ := os.OpenFile("/home/newmanh/go/src/decwars/syslog.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
	// uncomment the next line to get a log file
	// syscall.Dup2(int(logFile.Fd()), 1)

	syscall.Dup2(int(logFile.Fd()), 2)
	// fmt.Println("Starting Decwars...")

	//
	// Configure logger to write to the syslog. You could do this in init(), too.
	//
	f, err := os.OpenFile("/home/newmanh/go/src/decwars/decwars.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening log file: %v", err)
	}

	// don't forget to close it
	defer f.Close()

	// assign it to the standard logger
	log.SetOutput(f)
	//
	// Open the dbs
	//
	usersdb, objectsdb, pointsdb, robotsdb, planetsdb, basesdb := openDBs()

	//
	// Setup random # generator
	//
	rand.Seed(time.Now().UnixNano())
	//
	// build the universe
	//
	buildObjects(objectsdb, pointsdb, usersdb)
	//
	// Start up the robots
	//
	go robots(objectsdb, pointsdb, usersdb, robotsdb)
	//
	// Start up the planets
	//
	go planets(objectsdb, pointsdb, usersdb, planetsdb)
	//
	// Start up the bases
	//
	go bases(objectsdb, pointsdb, usersdb, basesdb)
	//
	// Start up the game ender process
	//
	go gamender(objectsdb, pointsdb, usersdb)

	p := ":" + strconv.Itoa(Portnum)
	//	ln, err := net.Listen("tcp", ":1701")
	ln, err := net.Listen("tcp", p)

	checkErr(err)

	tx, _ := atomicBegin(objectsdb)
	Conmap = make(map[string]Constr)
	err = atomicCommit(tx)
	//	Locmap = make(map[Loc]Player)
	//
	// Open/create long term storage of users
	//
	// Create game data
	//
	for {
		conn, err := ln.Accept()
		checkErr(err)
		go handleConnection(conn, usersdb, objectsdb, pointsdb)
		// fmt.Println("handling a connection")
	}
}
