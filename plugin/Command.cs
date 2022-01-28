using Sandbox.Game;
using System.Net;
using System.Net.Http;
using Torch.Commands;
using Torch.Commands.Permissions;
using VRage.Game.ModAPI;

namespace StalkR.LinkSteamDiscord
{
    public class Command : CommandModule
    {
        private LinkSteamDiscordPlugin _plugin => (LinkSteamDiscordPlugin)Context.Plugin;

        [Command("link", "Link your Steam & Discord accounts")]
        [Permission(MyPromoteLevel.None)]
        public async void link()
        {
            IMyPlayer player = Context.Player;
            if (player == null)
            {
                Context.Respond("Command must be run by a player (i.e. not console)");
                return;
            }

            string backend = this._plugin.Config.Backend;
            HttpResponseMessage response;
            using (HttpClient client = new HttpClient())
            {
                response = await client.GetAsync($"{backend}/steam/check/{player.SteamUserId}");
            }
            if (response.StatusCode == HttpStatusCode.OK)
            {
                Context.Respond($"Already linked, thanks! If you wish to forget or update the link, go to {backend}");
                return;
            }
            MyVisualScriptLogicProvider.OpenSteamOverlay($"https://steamcommunity.com/linkfilter/?url={backend}", player.IdentityId);
            Context.Respond("Browser window opened, please follow the steps to link your accounts.");
        }
    }
}
