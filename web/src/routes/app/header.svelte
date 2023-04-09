<script>
	import Sugarcube from '../../components/Sugarcube.svelte';
	import FaSignOutAlt from 'svelte-icons/fa/FaSignOutAlt.svelte';
</script>

<div
	class="flex w-full items-center gap-4 border-b border-stone-900 px-4 py-2 dark:border-stone-400"
>
	<div class="w-8">
		<Sugarcube style="stroke-width:15px;" />
	</div>
	<p class="font-mono text-2xl italic">sugarcube</p>
	<div class="flex-grow" />
	<div class="w-6">
		<FaSignOutAlt />
	</div>
</div>
