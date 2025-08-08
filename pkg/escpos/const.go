package escpos

type UnderlineMode byte

const (
	UnderlineOff       UnderlineMode = 0x00
	Underline1DotThick UnderlineMode = 0x01
	Underline2DotThick UnderlineMode = 0x02
)

type ItalicsMode byte

const (
	ItalicsOff ItalicsMode = 0x00
	ItalicsOn  ItalicsMode = 0x01
)

type EmphasisMode byte

const (
	EmphasisOff EmphasisMode = 0x00
	EmphasisOn  EmphasisMode = 0x01
)

type CharacterFont byte

const (
	CharacterFontA CharacterFont = 0x00
	CharacterFontB CharacterFont = 0x01
)

type CharacterCodePage int

const (
	CharacterCodePageDefault    CharacterCodePage = 0
	CharacterCodePagePC437      CharacterCodePage = 0
	CharacterCodePageKatakana   CharacterCodePage = 1
	CharacterCodePagePC850      CharacterCodePage = 2
	CharacterCodePagePC860      CharacterCodePage = 3
	CharacterCodePagePC863      CharacterCodePage = 4
	CharacterCodePagePC865      CharacterCodePage = 5
	CharacterCodePageHiragana   CharacterCodePage = 6
	CharacterCodePagePC851      CharacterCodePage = 11
	CharacterCodePagePC853      CharacterCodePage = 12
	CharacterCodePagePC857      CharacterCodePage = 13
	CharacterCodePagePC737      CharacterCodePage = 14
	CharacterCodePageISO8859_7  CharacterCodePage = 15
	CharacterCodePageWPC1252    CharacterCodePage = 16
	CharacterCodePagePC866      CharacterCodePage = 17
	CharacterCodePagePC852      CharacterCodePage = 18
	CharacterCodePagePC858      CharacterCodePage = 19
	CharacterCodePageISO8859_2  CharacterCodePage = 39
	CharacterCodePageISO8859_15 CharacterCodePage = 40
	CharacterCodePageWPC1250    CharacterCodePage = 45
)

type CutMode byte

const (
	CutModeFull    CutMode = 0x00
	CutModePartial CutMode = 0x01
)
