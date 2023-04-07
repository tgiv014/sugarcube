<script>
	import Footer from './footer.svelte';
	import Header from './header.svelte';

	import { status } from '../../lib/stores';
	import { goto } from '$app/navigation';

	$: if ($status) {
		if (!$status.setupCompleted) {
			goto('/setup');
		} else if (!$status.sessionValid) {
			goto('/');
		}
	}
</script>

<div class="flex flex-col min-h-screen w-full">
	<!-- todo: hide header/nav -->
	<Header />
	<div class="flex-1">
		<slot />
	</div>
	<Footer />
</div>
