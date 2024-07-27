import { Helmet } from 'react-helmet-async';
import { Box } from '@mui/material';
import Footer from 'src/components/footer';
import { PrivacyView } from 'src/sections/privacy/view';

export default function PrivacyPage() {
    return (
        <>
            <Helmet>
                <title> Privacy Policy | GooGL-Bye </title>
            </Helmet>

            <PrivacyView />
            <Footer />
        </>
    );
}
