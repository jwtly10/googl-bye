import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';
import Link from '@mui/material/Link';
import Divider from '@mui/material/Divider';
import { RouterLink } from 'src/routes/components';
import Logo from 'src/components/logo';

export default function TermsView() {
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
                        Terms and Conditions
                    </Typography>

                    <Typography variant="body1" sx={{ mb: 3 }}>
                        Last updated: July 27, 2024
                    </Typography>

                    <Typography variant="body1" paragraph>
                        Please read these Terms and Conditions carefully before using our website or
                        services.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        1. Acceptance of Terms
                    </Typography>
                    <Typography variant="body1" paragraph>
                        By accessing or using our GooGL-Bye service (hereinafter referred to as "the
                        Service"), you agree to be bound by these Terms and Conditions. If you
                        disagree with any part of these terms, you may not access the Service.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        2. Description of Service
                    </Typography>
                    <Typography variant="body1" paragraph>
                        The Service allows users to search for GitHub usernames and scan associated
                        repositories for goo.gl links. The Service also provides an option to raise
                        issues in the scanned repositories with the results of the parsing.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        3. Use of the Service
                    </Typography>
                    <Typography variant="body1" paragraph>
                        3.1 You agree to use the Service only for lawful purposes and in accordance
                        with these Terms.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        3.2 You must not use the Service in any way that causes, or may cause,
                        damage to the Service or impairment of the availability or accessibility of
                        the Service.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        3.3 You must not use the Service to copy, store, host, transmit, send, use,
                        publish or distribute any material which consists of (or is linked to) any
                        spyware, computer virus, Trojan horse, worm, keystroke logger, rootkit or
                        other malicious computer software.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        3.4 You must not conduct any systematic or automated data collection
                        activities on or in relation to the Service without our express written
                        consent.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        3.5 You must not use the Service for any purposes related to marketing
                        without our express written consent.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        4. Intellectual Property
                    </Typography>
                    <Typography variant="body1" paragraph>
                        4.1 All content on this website and the Service is the property of our
                        company and is protected by copyright laws.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        4.2 You may view, download for caching purposes only, and print pages from
                        the website for your own personal use, subject to the restrictions set out
                        below and elsewhere in these Terms and Conditions.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        5. Data Storage and Privacy
                    </Typography>
                    <Typography variant="body1" paragraph>
                        5.1 We store repository details such as repository name, stars, author, and
                        parsed links including their outbound destinations.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        5.2 By using the Service, you consent to this data collection and storage.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        5.3 We respect your privacy and handle your data in accordance with our
                        Privacy Policy.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        6. User Responsibilities
                    </Typography>
                    <Typography variant="body1" paragraph>
                        6.1 You are responsible for maintaining the confidentiality of your GitHub
                        credentials and any actions taken under your account.
                    </Typography>
                    <Typography variant="body1" paragraph>
                        6.2 You agree not to misuse the Service or help anyone else do so.
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        7. Limitation of Liability
                    </Typography>
                    <Typography variant="body1" paragraph>
                        7.1 To the fullest extent permitted by applicable law, we shall not be
                        liable for any indirect, incidental, special, consequential or punitive
                        damages, or any loss of profits or revenues, whether incurred directly or
                        indirectly, or any loss of data, use, goodwill, or other intangible losses
                        resulting from:
                    </Typography>
                    <Typography variant="body1" component="div">
                        <ul>
                            <li>your use or inability to use the Service;</li>
                            <li>
                                any unauthorized access to or use of our servers and/or any personal
                                information stored therein;
                            </li>
                            <li>
                                any interruption or cessation of transmission to or from the
                                Service;
                            </li>
                            <li>
                                any bugs, viruses, Trojan horses, or the like that may be
                                transmitted to or through the Service by any third party.
                            </li>
                        </ul>
                    </Typography>

                    <Typography variant="h5" sx={{ mt: 4, mb: 2 }}>
                        8. Changes to Terms
                    </Typography>
                    <Typography variant="body1" paragraph>
                        We reserve the right to modify these Terms at any time. We will always post
                        the most current version on our site. By continuing to use the Service after
                        changes have been made, you agree to be bound by the revised terms.
                    </Typography>

                    <Divider sx={{ my: 4 }} />

                    <Typography variant="body2" sx={{ mt: 2, color: 'text.secondary' }}>
                        If you have any questions about these Terms and Conditions, please contact
                        me on twitter at <Link href="https://x.com/jwtly10">jwtly10</Link>.
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
