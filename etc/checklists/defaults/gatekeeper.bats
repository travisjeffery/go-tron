load test_helper

@test "Gatekeeper enabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run spctl --status
  [ $status -eq 0 ] && [ $output = "assessments enabled" ]
}
