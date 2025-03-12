package sounds

import blueprintclient "github.com/TheBitDrifter/blueprint/client"

var Run = blueprintclient.SoundConfig{
	Path: "run.wav",
	// We could increase the count for two players, but it sounds too clustered in my opinion
	AudioPlayerCount: 1,
}

var Music = blueprintclient.SoundConfig{
	Path:             "fantasy_music.wav",
	AudioPlayerCount: 1,
}
