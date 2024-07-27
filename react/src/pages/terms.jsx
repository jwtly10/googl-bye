import { Helmet } from 'react-helmet-async';

import { TermsView } from 'src/sections/terms/view';

import Footer from 'src/components/footer';

// ----------------------------------------------------------------------

export default function TermsPage() {
    return (
        <>
            <Helmet>
                <title> Terms & Conditions | GooGL-Bye </title>
            </Helmet>

            <TermsView />
            <Footer />
        </>
    );
}
