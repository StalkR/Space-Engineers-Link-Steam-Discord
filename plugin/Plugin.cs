using System.IO;
using System.Windows.Controls;
using Torch;
using Torch.API;
using Torch.API.Plugins;
using VRage.Utils;

namespace StalkR.LinkSteamDiscord
{
    public class Plugin : TorchPluginBase, IWpfPlugin
    {
        private Persistent<Config> _config;
        public Config Config => _config?.Data;
        public void Save() => _config?.Save();

        private UserControl _control;
        public UserControl GetControl() => _control ?? (_control = new ConfigUI(this));

        public override void Init(ITorchBase torch)
        {
            base.Init(torch);

            string path = Path.Combine(StoragePath, "LinkSteamDiscord.cfg");
            MyLog.Default.WriteLine($"LinkSteamDiscord: loading config from {path}");
            _config = Persistent<Config>.Load(path);
        }
    }
}
