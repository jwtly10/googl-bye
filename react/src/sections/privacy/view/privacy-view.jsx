import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';
import Divider from '@mui/material/Divider';
import { RouterLink } from 'src/routes/components';
import Logo from 'src/components/logo';

export default function PrivacyView() {
    const renderHeader = (
        <Box
            component="header"
            sx={{
                top: 0,
                left: 0,
                width: 1,
                lineHeight: 0,
                position: 'fixed',
                p: (theme) => ({ xs: theme.spacing(3, 3, 0), sm: theme.spacing(5, 5, 0) }),
            }}
        >
            <Logo />
        </Box>
    );

    return (
        <>
            {renderHeader}
            <Container>
                <Box
                    sx={{
                        py: 12,
                        maxWidth: 800,
                        mx: 'auto',
                        minHeight: '100vh',
                        textAlign: 'left',
                        display: 'flex',
                        flexDirection: 'column',
                    }}
                >
                    <Typography variant="h3" sx={{ mb: 4 }}>
                        Privacy Policy
                    </Typography>

                    <Typography variant="body1" sx={{ mb: 3 }}>
                        Last updated: July 27, 2024
                    </Typography>

                    <Typography variant="body1" paragraph>
                        Your privacy is important to us. This Privacy Policy explains how we
                        collect, use, disclose, and safeguard your information when you use our
                        website or services.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        1. Introduction
                    </Typography>
                    <Typography variant="body1" paragraph>
                        This Privacy Policy explains how we collect, use, process, and disclose your
                        information in connection with our GooGL-Bye service.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        2. Information We Collect
                    </Typography>
                    <Typography variant="body1" paragraph>
                        2.1 GitHub Information: When you use our service, we collect and store
                        information related to the GitHub repositories you search, including:
                    </Typography>
                    <Typography variant="body1" component="div">
                        <ul>
                            <li>Repository names</li>
                            <li>Star counts</li>
                            <li>Fork counts</li>
                            <li>Author information</li>
                            <li>goo.gl links found in the repositories</li>
                            <li>Outbound links associated with the goo.gl links</li>
                        </ul>
                    </Typography>
                    <Typography variant="body1" paragraph>
                        2.2 Usage Information: We may collect information about how you use our
                        service, including your search queries and interaction with the features we
                        provide.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        3. How We Use Your Information
                    </Typography>
                    <Typography variant="body1" paragraph>
                        We use the collected information to:
                    </Typography>
                    <Typography variant="body1" component="div">
                        <ul>
                            <li>Provide and maintain our service</li>
                            <li>Improve and optimize our service</li>
                            <li>Respond to your comments or questions</li>
                            <li>Provide support for the service</li>
                        </ul>
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        4. Information Sharing and Disclosure
                    </Typography>
                    <Typography variant="body1" paragraph>
                        We do not share or sell your personal information to third parties.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        5. Data Retention
                    </Typography>
                    <Typography variant="body1" paragraph>
                        We retain the collected information for as long as your account is active or
                        as needed to provide you with our services. We will delete this information
                        upon your request, except where we are required to retain it by law.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        6. Security
                    </Typography>
                    <Typography variant="body1" paragraph>
                        We implement reasonable security measures to protect the security of your
                        information. However, please be aware that no security measures are perfect
                        or impenetrable.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        7. Your Rights
                    </Typography>
                    <Typography variant="body1" paragraph>
                        You have the right to access, update, or delete your information. If you
                        wish to exercise these rights, please contact us using the information
                        provided below.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        8. Changes to This Privacy Policy
                    </Typography>
                    <Typography variant="body1" paragraph>
                        We may update our Privacy Policy from time to time. We will notify you of
                        any changes by posting the new Privacy Policy on this page and updating the
                        "Last Updated" date.
                    </Typography>

                    <Divider sx={{ my: 4 }} />

                    <Typography variant="body2" sx={{ mt: 2, color: 'text.secondary' }}>
                        If you have any questions about this Privacy Policy, please contact me on
                        twitter at <Link href="https://x.com/jwtly10">jwtly10</Link>.
                    </Typography>

                    <Button
                        href="/"
                        size="large"
                        variant="contained"
                        component={RouterLink}
                        sx={{ mt: 4, alignSelf: 'flex-start' }}
                    >
                        Back to Home
                    </Button>
                </Box>
            </Container>
        </>
    );
}
