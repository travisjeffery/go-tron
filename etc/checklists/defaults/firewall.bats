load test_helper

@test "Firewall enabled" {  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read /Library/Preferences/com.apple.alf globalstate
  [ $status -eq 0 ] && [ $output -eq 1 ]
}
