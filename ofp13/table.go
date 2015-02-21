package ofp13

import (
	"io"

	"github.com/netrack/openflow/encoding/binary"
)

const (
	MAX_TABLE_NAME_LEN = 32

	// Last usable table number
	TT_MAX Table = 0xfe

	// Fake tables
	TT_ALL Table = 0xff
)

type Table uint8

const (
	TC_DEPRECATED_MASK TableConfig = 3
)

type TableConfig uint32

// Configure/Modify behavior of a flow table
type TableMod struct {
	// The table_id chooses the table to which the configuration
	// change should be applied. If the TableID is OFPTT_ALL,
	// the configuration is applied to all tables in the switch.
	TableID Table
	// The config field is a bitmap that is provided for backward
	// compatibility with earlier version of the specification,
	// it is reserved for future use.
	Config TableConfig
}

func (t *TableMod) ReadFrom(r io.Reader) (int64, error) {
	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&t.TableID,
		&pad3{},
		&t.Config,
	})
}

func (t *TableMod) WriteTo(w io.Writer) (int64, error) {
	return binary.WriteSlice(w, binary.BigEndian, []interface{}{
		t.TableID,
		pad3{},
		t.Config,
	})
}

// Information about tables is requested with the MP_TABLE multipart
// request type. The request does not contain any data in the body.
// The body of the reply consists of an array of the TableStats
type TableStats struct {
	// Identifier of table. Lower numbered tables are consulted first
	TableID Table
	// Number of active entries
	ActiveCount uint32
	// Number of packets looked up in table
	LookupCount uint64
	// Number of packets that hit table
	MatchedCount uint64
}

func (t *TableStats) ReadFrom(r io.Reader) (int64, error) {
	return binary.ReadSlice(r, binary.BigEndian, []interface{}{
		&t.TableID,
		&pad3{},
		&t.ActiveCount,
		&t.LookupCount,
		&t.MatchedCount,
	})
}

type TableFeatures struct {
	Length  uint16
	TableID Table
	_       pad5
	name    [MAX_TABLE_NAME_LEN]byte

	MetadataMatch uint64
	MetadataWrite uint64
	Config        TableConfig

	MaxEntries uint32
	Properties []TableFeaturePropHeader
}

const (
	TFPT_INSTRUCTIONS TableFeaturePropType = iota
	TFPT_INSTRUCTIONS_MISS
	TFPT_NEXT_TABLES
	TFPT_NEXT_TABLES_MISS
	TFPT_WRITE_ACTIONS
	TFPT_WRITE_ACTIONS_MISS
	TFPT_APPLY_ACTIONS
	TFPT_APPLY_ACTIONS_MISS
	TFPT_MATCH
	TFPT_WILDCARDS
	TFPT_WRITE_SETFIELD
	TFPT_WRITE_SETFIELD_MISS
	TFPT_APPLY_SETFIELD
	TFPT_APPLY_SETFIELD_MISS
	TFPT_EXPERIMENTER      TableFeaturePropType = 0xfffe
	TFPT_EXPERIMENTER_MISS TableFeaturePropType = 0xffff
)

type TableFeaturePropType uint16

type TableFeaturePropHeader struct {
	Type   TableFeaturePropType
	Length uint16
}

type TableFeaturePropInstructions struct {
	Type   TableFeaturePropType
	Length uint16
	//TODO: InstructionID  []Instruction
}

type TableFeaturePropNextTables struct {
	Type        TableFeaturePropType
	Length      uint16
	NextTableID []Table
}

type TableFeaturePropActions struct {
	Type     TableFeaturePropType
	Length   uint16
	ActionID []interface{}
}

type TableFeaturePropOXM struct {
	Type   TableFeaturePropType
	Length uint16
	OXMID  []OXM
}

type TableFeaturePropExperimenter struct {
	Type             TableFeaturePropType
	Length           uint16
	Experimenter     uint32
	ExpType          uint32
	ExperimenterData []uint32
}