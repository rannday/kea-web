package utils

import (
	"flag"
)

func ParseCLI(env *Env) {
  if env == nil {
    return
  }

  flag.StringVar(
    &env.PORT,
    "port",
    env.PORT,
    "HTTP listen port",
  )

  flag.StringVar(
    &env.STATIC_DIR,
    "static-dir",
    env.STATIC_DIR,
    "Static directory override",
  )
  flag.Parse()
}
