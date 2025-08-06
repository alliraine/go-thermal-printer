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
	CharacterCodePageDefault           CharacterCodePage = 0
	CharacterCodePageCP437             CharacterCodePage = 3
	CharacterCodePageCP808             CharacterCodePage = 17
	CharacterCodePageGeorgianMkhedruli CharacterCodePage = 18
)

type CutMode byte

const (
	CutModeFull    CutMode = 0x00
	CutModePartial CutMode = 0x01
)
