package runner

import (
  "fmt"
  "github.com/pilu/config"
  "os"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)

const (
  envSettingsPrefix   = "RUNNER_"
  mainSettingsSection = "Settings"
)

var settings = map[string]string{
  "config_path":       "./runner.conf",
  "root":              ".",
  "tmp_path":          "./tmp",
  "build_name":        "runner-build",
  "build_log":         "runner-build-errors.log",
  "valid_ext":         ".go, .tpl, .tmpl, .html",
  "build_delay":       "600",
  "colors":            "1",
  "log_color_main":    "cyan",
  "log_color_build":   "yellow",
  "log_color_runner":  "green",
  "log_color_watcher": "magenta",
  "log_color_app":     "",
  "exclude_dir":       ""
}

var excludedDirs map[string]bool

var colors = map[string]string{
  "reset":          "0",
  "black":          "30",
  "red":            "31",
  "green":          "32",
  "yellow":         "33",
  "blue":           "34",
  "magenta":        "35",
  "cyan":           "36",
  "white":          "37",
  "bold_black":     "30;1",
  "bold_red":       "31;1",
  "bold_green":     "32;1",
  "bold_yellow":    "33;1",
  "bold_blue":      "34;1",
  "bold_magenta":   "35;1",
  "bold_cyan":      "36;1",
  "bold_white":     "37;1",
  "bright_black":   "30;2",
  "bright_red":     "31;2",
  "bright_green":   "32;2",
  "bright_yellow":  "33;2",
  "bright_blue":    "34;2",
  "bright_magenta": "35;2",
  "bright_cyan":    "36;2",
  "bright_white":   "37;2",
}

func logColor(logName string) string {
  settingsKey := fmt.Sprintf("log_color_%s", logName)
  colorName := settings[settingsKey]

  return colors[colorName]
}

func loadEnvSettings() {
  for key, _ := range settings {
    envKey := fmt.Sprintf("%s%s", envSettingsPrefix, strings.ToUpper(key))
    if value := os.Getenv(envKey); value != "" {
      settings[key] = value
    }
  }
}

func loadRunnerConfigSettings() {
  if _, err := os.Stat(configPath()); err != nil {
    return
  }

  logger.Printf("Loading settings from %s", configPath())
  sections, err := config.ParseFile(configPath(), mainSettingsSection)
  if err != nil {
    return
  }

  for key, value := range sections[mainSettingsSection] {
    settings[key] = value
  }
}

func initExcludedDirs() {
  excludedDirs = make(map[string]bool)
  for _, dir := range strings.Split(settings["exclude_dir"], ",") {
    dp, err := filepath.Abs(dir)
    if err != nil {
      continue
    }
    excludedDirs[dp] = true
    fmt.Printf("Directory %s is not watched\n", dp)
  }
}

func initSettings() {
  loadEnvSettings()
  loadRunnerConfigSettings()
  initExcludedDirs()
}

func getenv(key, defaultValue string) string {
  if value := os.Getenv(key); value != "" {
    return value
  }

  return defaultValue
}

func root() string {
  return settings["root"]
}

func tmpPath() string {
  return settings["tmp_path"]
}

func buildName() string {
  return settings["build_name"]
}
func buildPath() string {
  return filepath.Join(tmpPath(), buildName())
}

func buildErrorsFileName() string {
  return settings["build_log"]
}

func buildErrorsFilePath() string {
  return filepath.Join(tmpPath(), buildErrorsFileName())
}

func configPath() string {
  return settings["config_path"]
}

func buildDelay() time.Duration {
  value, _ := strconv.Atoi(settings["build_delay"])

  return time.Duration(value)
}
