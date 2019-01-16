//go:generate $PWD/bin/pbgraphserialization-gen --targetPackagePath github.com/loomnetwork/gamechain/battleground/game/ --targetPackageName game --protoPackageName zb --outputPath types_serialization.go

package game

import _ "github.com/loomnetwork/gamechain/types/zb"
import _ "github.com/loomnetwork/gamechain/library/pbgraphserialization"

1