load test_helper

@test "Password policy requires alpha, numeric, and 8 chars minimum" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run pwpolicy getglobalpolicy
  echo "$output" | grep "requiresAlpha=1"
  echo "$output" | grep "requiresNumeric=1"
  (( $(echo "$output" | grep -E -o "minChars=(\d+)" | cut -d= -f 2) >= 8 ))
}

@test "Fast user switching disabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read /Library/Preferences/.GlobalPreferences MultipleSessionEnabled
  [ $status -ne 0 ] || [ "$output" -eq 0 ]
}

@test "Root user login disabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run dscl . -read /Users/root AuthenticationAuthority
  [ $status -eq 0 ] || [ $output = "No such key: AuthenticationAuthority" ]
}

@test "Automatic login disabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read /Library/Preferences/com.apple.loginwindow autoLoginUser
  [ $status -eq 1 ]
}

@test "Password hint disabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read /Library/Preferences/com.apple.loginwindow RetriesUntilHint
  [ $status -ne 0 ] || [ "$output" -eq 0 ]
}

@test "Guest login disabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read /Library/Preferences/com.apple.loginwindow GuestEnabled
  [ $status -ne 0 ] || [ "$output" -eq 0 ]
}

@test "Requires password 5 minutes (or less) after sleep or screen saver begins" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read com.apple.screensaver askForPassword
  [ $status -eq 0 ] && [ "$output" -eq 1 ]
  run defaults read com.apple.screensaver askForPasswordDelay
  [ $status -ne 0 ] || [ "$output" -le 300 ]
}
