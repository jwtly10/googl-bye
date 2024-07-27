import PropTypes from 'prop-types';
import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import { useTheme } from '@mui/material/styles';
import IconButton from '@mui/material/IconButton';
import { useResponsive } from 'src/hooks/use-responsive';
import { bgBlur } from 'src/theme/css';
import Iconify from 'src/components/iconify';
import Searchbar from './common/searchbar';
import { NAV, HEADER } from './config-layout';
import AccountPopover from './common/account-popover';
import LanguagePopover from './common/language-popover';
import NotificationsPopover from './common/notifications-popover';

// ----------------------------------------------------------------------

export default function Header({ onOpenNav }) {
    const theme = useTheme();
    const lgUp = useResponsive('up', 'lg');

    const renderContent = (
        <>
            <Box
                sx={{
                    display: 'flex',
                    alignItems: 'center',
                    flexGrow: 1,
                }}
            >
                <Box
                    component="img"
                    src="/logo.png"
                    alt="G Logo"
                    sx={{
                        height: 30,
                        marginRight: '0.2rem',
                    }}
                />
                <Typography
                    variant="h4"
                    component="div"
                    sx={{
                        fontWeight: 'bold',
                        fontFamily: 'Poppins, sans-serif',
                        letterSpacing: 1,
                        color: theme.palette.primary.main,
                        '& span': {
                            color: theme.palette.secondary.main,
                        },
                    }}
                >
                    ooGL-<span>Bye</span>
                </Typography>
            </Box>
            <Stack direction="row" alignItems="center" spacing={1}>
                <AccountPopover />
            </Stack>
        </>
    );

    return (
        <AppBar
            sx={{
                boxShadow: theme.shadows.z8,
                height: HEADER.H_MOBILE,
                zIndex: theme.zIndex.appBar + 1,
                ...bgBlur({
                    color: theme.palette.background.default,
                }),
                transition: theme.transitions.create(['height'], {
                    duration: theme.transitions.duration.shorter,
                }),
            }}
        >
            <Toolbar
                sx={{
                    height: 1,
                    px: { lg: 5 },
                }}
            >
                {renderContent}
            </Toolbar>
        </AppBar>
    );
}

Header.propTypes = {
    onOpenNav: PropTypes.func,
};
