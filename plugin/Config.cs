using Torch;

namespace StalkR.LinkSteamDiscord
{
    public class Config : ViewModel
    {
        private string _backend = "https://link-steam-discord.stalkr.net/";
        public string Backend { get => _backend.TrimEnd('/'); set => SetValue(ref _backend, value); }
    }
}