using NLog;
using System.IO;
using System.Windows.Controls;
using Torch;
using Torch.API;
using Torch.API.Plugins;

namespace StalkR.LinkSteamDiscord
{
    public class LinkSteamDiscordPlugin : TorchPluginBase, IWpfPlugin
    {
        private static readonly Logger Log = LogManager.GetCurrentClassLogger();

        private Persistent<Config> _config;
        public Config Config => _config?.Data;
        public void Save() => _config?.Save();

        private UserControl _control;
        public UserControl GetControl() => _control ?? (_control = new ConfigUI(this));

        public override void Init(ITorchBase torch)
        {
            base.Init(torch);

            string path = Path.Combine(StoragePath, "LinkSteamDiscord.cfg");
            Log.Info($"Attempting to load config from {path}");
            _config = Persistent<Config>.Load(path);
        }
    }
}
