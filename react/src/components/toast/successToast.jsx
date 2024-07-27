import { Snackbar, Alert, Link } from '@mui/material';
import { useState, useEffect } from 'react';

export default function SuccessToast({ message, url, open, onClose }) {
    const [isOpen, setIsOpen] = useState(open);

    useEffect(() => {
        setIsOpen(open);
    }, [open]);

    const handleClose = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        setIsOpen(false);
        if (onClose) {
            onClose();
        }
    };

    const renderMessage = () => {
        if (url) {
            console.log('rendering url');
            return (
                <>
                    {message}{' '}
                    <Link
                        href={url}
                        target="_blank"
                        rel="noopener noreferrer"
                        color="inherit"
                        underline="always"
                    >
                        {url}
                    </Link>
                </>
            );
        }
        return message;
    };

    return (
        <Snackbar
            open={isOpen}
            autoHideDuration={6000}
            onClose={handleClose}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            sx={{ zIndex: (theme) => theme.zIndex.drawer + 2 }}
        >
            <Alert
                onClose={handleClose}
                severity="success"
                variant="filled"
                sx={{
                    width: '100%',
                    bgcolor: 'success.dark',
                    color: 'white',
                    '& .MuiAlert-icon': {
                        color: 'white',
                    },
                    '& a': {
                        color: 'white',
                        fontWeight: 'bold',
                    },
                }}
            >
                {renderMessage()}
            </Alert>
        </Snackbar>
    );
}
