load test_helper

@test "FileVault enabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run fdesetup status
  [ $status -eq 0 ] && [ "$output" == "FileVault is On." ]
}
