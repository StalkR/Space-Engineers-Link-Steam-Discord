using System.Windows;
using System.Windows.Controls;

namespace StalkR.LinkSteamDiscord
{

    public partial class ConfigUI : UserControl
    {
        private Plugin Plugin { get; }

        public ConfigUI()
        {
            InitializeComponent();
        }

        public ConfigUI(Plugin plugin) : this()
        {
            Plugin = plugin;
            DataContext = plugin.Config;
        }

        private void SaveConfig_OnClick(object sender, RoutedEventArgs e)
        {
            Plugin.Save();
        }
    }
}
