<script>
    import SpectatorView from "@/views/SpectatorView.svelte";
    import { username } from "@/stores";
    import { Team } from "@/utils";

    username.set("Adrien")

</script>

<SpectatorView team={Team.BLUE} cards={["Adrien"]}/>
