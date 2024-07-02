<script>
    import { Route } from "tinro";
    import _instance from "./stream";
    import Dialer, { ArgsConstructor, FetchActive } from "./utils";

    import Dashboard from "./views/Dashboard.svelte";
    import Inventories from "./views/Inventories.svelte";
    import Sales from "./views/Sales.svelte";
    import AdminSignup from "./views/AdminSignup.svelte";
    import ResetPassword from "./views/ResetPassword.svelte";
    import LoginUser from "./views/LoginUser.svelte";
    import Reports from "./views/Reports.svelte";
    import StaffManagement from "./views/StaffManagement.svelte";
    import CustomerRelations from "./views/CustomerRelations.svelte";
    import CreateNewTab from "./views/CreateNewTab.svelte";
    import Settings from "./views/Settings.svelte";

	console.log(window.location.pathname);
</script>

<main>
	{#await Dialer(Exists, ArgsConstructor(["admin", [["name", "admin"]]])) then res}
		{#if !res}
    		<Route path="/*"><AdminSignup signUpMessage="Set up your administrator account" /></Route>
		{:else}
			{#if !FetchActive}
				<Route path="/*"><LoginUser /></Route>
			{:else}
    		<!--<Route path="/"><Dashboard /></Route>-->
			<Route path="/login"><LoginUser /></Route>
			<Route path="/admin_signup"><AdminSignup signUpMessage="An admin account already exists" /></Route>
			<Route path="/reset_password"><ResetPassword /></Route>
			<Route path="/"><Sales /></Route>
			<Route path="/reports"><Reports /></Route>
			<Route path="/settings/:settings"><Settings /></Route>
    		<Route path="/inventories/:tabData"><Inventories /></Route>
			<Route path="/staff/:staff"><StaffManagement /></Route>
			<Route path="/customers/:customers"><CustomerRelations /></Route>
			<Route path="/tabs/:tabs"><CreateNewTab /></Route>
			{/if}
		{/if}
	{/await}
</main>

<style>
</style>
