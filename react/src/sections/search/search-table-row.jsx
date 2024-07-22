import { useState } from 'react';
import PropTypes from 'prop-types';
import Stack from '@mui/material/Stack';
import Avatar from '@mui/material/Avatar';
import Popover from '@mui/material/Popover';
import TableRow from '@mui/material/TableRow';
import Checkbox from '@mui/material/Checkbox';
import MenuItem from '@mui/material/MenuItem';
import TableCell from '@mui/material/TableCell';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import Label from 'src/components/label';
import Iconify from 'src/components/iconify';

// ----------------------------------------------------------------------

export default function RepoTableRow({
    selected,
    name,
    author,
    language,
    stars,
    forks,
    parseStatus,
    handleClick,
}) {
    const [open, setOpen] = useState(null);

    const handleOpenMenu = (event) => {
        setOpen(event.currentTarget);
    };

    const handleCloseMenu = () => {
        setOpen(null);
    };

    return (
        <>
            <TableRow hover tabIndex={-1} role="checkbox" selected={selected}>
                <TableCell padding="checkbox">
                    <Checkbox disableRipple checked={selected} onChange={handleClick} />
                </TableCell>
                <TableCell>
                    <Typography variant="subtitle2" noWrap>
                        {name}
                    </Typography>
                </TableCell>
                <TableCell>{author}</TableCell>
                <TableCell>{language}</TableCell>
                <TableCell align="center">{stars}</TableCell>
                <TableCell align="center">{forks}</TableCell>
                <TableCell>
                    <Label
                        color={
                            (parseStatus === 'PENDING' && 'warning') ||
                            (parseStatus === 'PROCESSING' && 'info') ||
                            (parseStatus === 'DONE' && 'success') ||
                            (parseStatus === 'ERROR' && 'error') ||
                            'default'
                        }
                    >
                        {parseStatus}
                    </Label>
                </TableCell>
                <TableCell align="right">
                    <IconButton onClick={handleOpenMenu}>
                        <Iconify icon="eva:more-vertical-fill" />
                    </IconButton>
                </TableCell>
            </TableRow>
            <Popover
                open={!!open}
                anchorEl={open}
                onClose={handleCloseMenu}
                anchorOrigin={{ vertical: 'top', horizontal: 'left' }}
                transformOrigin={{ vertical: 'top', horizontal: 'right' }}
                PaperProps={{
                    sx: { width: 140 },
                }}
            >
                <MenuItem onClick={handleCloseMenu}>
                    <Iconify icon="eva:edit-fill" sx={{ mr: 2 }} />
                    Edit
                </MenuItem>
                <MenuItem onClick={handleCloseMenu} sx={{ color: 'error.main' }}>
                    <Iconify icon="eva:trash-2-outline" sx={{ mr: 2 }} />
                    Delete
                </MenuItem>
            </Popover>
        </>
    );
}

RepoTableRow.propTypes = {
    avatarUrl: PropTypes.string,
    author: PropTypes.string,
    handleClick: PropTypes.func,
    language: PropTypes.string,
    name: PropTypes.string,
    stars: PropTypes.number,
    forks: PropTypes.number,
    selected: PropTypes.bool,
    parseStatus: PropTypes.string,
};
