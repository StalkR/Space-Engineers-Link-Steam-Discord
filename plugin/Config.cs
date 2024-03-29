using Torch;

namespace StalkR.LinkSteamDiscord
{
    public class Config : ViewModel
    {
        private string _backend = "https://example.com/";
        public string Backend { get => _backend.TrimEnd('/'); set => SetValue(ref _backend, value); }
    }
}
