<script>
	import Sugarcube from '../../components/Sugarcube.svelte';
	import FaSignOutAlt from 'svelte-icons/fa/FaSignOutAlt.svelte';
	import MdSettings from 'svelte-icons/md/MdSettings.svelte';
	import { logout } from '$lib/auth';
</script>

<div
	class="flex w-full items-center gap-4 border-b border-stone-900 px-4 py-2 dark:border-stone-400"
>
	<a href="/app" class="flex gap-4">
		<div class="w-8">
			<Sugarcube style="stroke-width:15px;" />
		</div>
		<p class="font-mono text-2xl italic">sugarcube</p>
	</a>

	<div class="flex-grow" />
	<a href="/app/settings" class="w-6">
		<MdSettings />
	</a>
	<button
		on:click={() => {
			logout().then(() => {
				console.log('logout');
			});
		}}
		class="w-6"
	>
		<FaSignOutAlt />
	</button>
</div>
