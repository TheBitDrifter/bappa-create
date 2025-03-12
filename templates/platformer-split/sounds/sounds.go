package sounds

import blueprintclient "github.com/TheBitDrifter/blueprint/client"

var Run = blueprintclient.SoundConfig{
	Path:             "run.wav",
	AudioPlayerCount: 1, // Too much running sounds clustered, keep at one
}

var Jump = blueprintclient.SoundConfig{
	Path:             "jump.wav",
	AudioPlayerCount: 2, // Increase to player count
}

var Land = blueprintclient.SoundConfig{
	Path:             "land.wav",
	AudioPlayerCount: 2, // Increase to player count
}

var Music = blueprintclient.SoundConfig{
	Path:             "music.wav",
	AudioPlayerCount: 1,
}
