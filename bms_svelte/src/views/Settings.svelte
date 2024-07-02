<script>
  import Grid from "../components/Grid.svelte";
import Nav from "../components/Nav.svelte";

let backup = () => {
	
}, resetPass = () => window.location.assign("/reset_password")
</script>

<main>
<Nav />

<!--<div class="card">
    <div class="setting">
        <p class="setting-name"><span class="material-symbols-outlined">backup</span></p>
        <p class="current-setting">Backup data</p>
    </div>
</div>-->

<div class="card" on:click={resetPass}>
    <div class="setting">
        <p class="setting-name"><span class="material-symbols-outlined">password</span></p>
        <p class="current-setting">Change admin password</p>
    </div>
</div>

<h4 style="text-align: center;">Mpesa till number . Qr code scanner . Receipt printer | Settings</h4>
<hr>
<Grid tabSplit="10"/>
</main>

<style>
    .card {
        padding: 10px;
        border: 1px solid #ccc;
        border-radius: 8px;
        cursor: pointer;
        transition: transform 0.3s ease-in-out;
		margin: 20px 200px;
    }
	
	.card:active {
		transform: translateY(6px);
	}

    .setting {
        display: flex;
        justify-content: space-between;
    }

    .setting-name {
        font-size: 16px;
        font-weight: bold;
        color: #333;
        margin-right: 10px;
    }

    .current-setting {
        font-size: 16px;
        color: #555;
		font-weight: bolder;
    }
</style>
