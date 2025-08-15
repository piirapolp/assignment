INSERT INTO `health` (`k`, `v`) VALUES ('mock_data', 'complete')
    ON DUPLICATE KEY UPDATE v = 'complete';