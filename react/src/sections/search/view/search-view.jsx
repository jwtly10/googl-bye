import { useState } from 'react';

import {
    Stack,
    Table,
    Container,
    TableBody,
    TableContainer,
    TablePagination,
    Button,
    Typography,
    Card,
} from '@mui/material';

import { users } from 'src/_mock/user';

import Iconify from 'src/components/iconify';
import Scrollbar from 'src/components/scrollbar';

import SearchForm from 'src/components/search/searchForm';
import ErrorToast from 'src/components/toast/errorToast';

import TableNoData from '../table-no-data';
import UserTableRow from '../search-table-row';
import UserTableHead from '../search-table-head';
import TableEmptyRows from '../table-empty-rows';
import UserTableToolbar from '../search-table-toolbar';
import { emptyRows, applyFilter, getComparator } from '../utils';

// ----------------------------------------------------------------------

export default function SearchPage() {
    const [page, setPage] = useState(0);

    const [order, setOrder] = useState('asc');

    const [selected, setSelected] = useState([]);

    const [orderBy, setOrderBy] = useState('name');

    const [filterName, setFilterName] = useState('');

    const [rowsPerPage, setRowsPerPage] = useState(5);

    const handleSort = (event, id) => {
        const isAsc = orderBy === id && order === 'asc';
        if (id !== '') {
            setOrder(isAsc ? 'desc' : 'asc');
            setOrderBy(id);
        }
    };

    const handleSelectAllClick = (event) => {
        if (event.target.checked) {
            const newSelecteds = users.map((n) => n.name);
            setSelected(newSelecteds);
            return;
        }
        setSelected([]);
    };

    const handleClick = (event, name) => {
        const selectedIndex = selected.indexOf(name);
        let newSelected = [];
        if (selectedIndex === -1) {
            newSelected = newSelected.concat(selected, name);
        } else if (selectedIndex === 0) {
            newSelected = newSelected.concat(selected.slice(1));
        } else if (selectedIndex === selected.length - 1) {
            newSelected = newSelected.concat(selected.slice(0, -1));
        } else if (selectedIndex > 0) {
            newSelected = newSelected.concat(
                selected.slice(0, selectedIndex),
                selected.slice(selectedIndex + 1)
            );
        }
        setSelected(newSelected);
    };

    const handleChangePage = (event, newPage) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event) => {
        setPage(0);
        setRowsPerPage(parseInt(event.target.value, 10));
    };

    const handleFilterByName = (event) => {
        setPage(0);
        setFilterName(event.target.value);
    };

    const dataFiltered = applyFilter({
        inputData: users,
        comparator: getComparator(order, orderBy),
        filterName,
    });

    const [isSearching, setIsSearching] = useState(false);
    const [errorToast, setErrorToast] = useState({ open: false, message: '' });

    const [searchParams, setSearchParams] = useState({
        name: '',
        query: '',
        opts: {
            sort: '',
            order: '',
        },
        startPage: 0,
        currentPage: 0,
        pagesToProcess: 1,
    });

    const validateForm = () => {
        if (!searchParams.name.trim()) {
            return 'Search Name is required';
        }

        if (!searchParams.query.trim()) {
            return 'Query is required';
        }

        if (searchParams.startPage < 0) {
            return 'Start Page must be at least 0';
        }

        if (searchParams.pagesToProcess < 1) {
            return 'Pages to Process must be at least 1';
        }

        if (searchParams.pagesToProcess > 20) {
            return 'Pages to Process is limited to 20';
        }

        return null;
    };

    const handleCloseErrorToast = () => {
        setErrorToast({ open: false, message: '' });
    };

    const handleSearch = () => {
        const error = validateForm();

        if (error) {
            setErrorToast({ open: true, message: error });
            return;
        }

        setIsSearching(true);
        console.log('Making github api search');
        console.log(searchParams);
        console.log('Starting search...');

        // Here is were we implement search
        setTimeout(() => {
            setIsSearching(false);
            console.log('Search completed!');
        }, 3000);
    };

    const notFound = !dataFiltered.length && !!filterName;

    return (
        <Container>
            <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
                <Typography variant="h4">Search</Typography>
            </Stack>

            <SearchForm
                setSearchParams={setSearchParams}
                handleSearch={handleSearch}
                searchParams={searchParams}
                isSearching={isSearching}
            />

            <Card>
                <UserTableToolbar
                    numSelected={selected.length}
                    filterName={filterName}
                    onFilterName={handleFilterByName}
                />

                <Scrollbar>
                    <TableContainer sx={{ overflow: 'unset' }}>
                        <Table sx={{ minWidth: 800 }}>
                            <UserTableHead
                                order={order}
                                orderBy={orderBy}
                                rowCount={users.length}
                                numSelected={selected.length}
                                onRequestSort={handleSort}
                                onSelectAllClick={handleSelectAllClick}
                                headLabel={[
                                    { id: 'name', label: 'Name' },
                                    { id: 'company', label: 'Company' },
                                    { id: 'role', label: 'Role' },
                                    { id: 'isVerified', label: 'Verified', align: 'center' },
                                    { id: 'status', label: 'Status' },
                                    { id: '' },
                                ]}
                            />
                            <TableBody>
                                {dataFiltered
                                    .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                                    .map((row) => (
                                        <UserTableRow
                                            key={row.id}
                                            name={row.name}
                                            role={row.role}
                                            status={row.status}
                                            company={row.company}
                                            avatarUrl={row.avatarUrl}
                                            isVerified={row.isVerified}
                                            selected={selected.indexOf(row.name) !== -1}
                                            handleClick={(event) => handleClick(event, row.name)}
                                        />
                                    ))}

                                <TableEmptyRows
                                    height={77}
                                    emptyRows={emptyRows(page, rowsPerPage, users.length)}
                                />

                                {notFound && <TableNoData query={filterName} />}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Scrollbar>

                <TablePagination
                    page={page}
                    component="div"
                    count={users.length}
                    rowsPerPage={rowsPerPage}
                    onPageChange={handleChangePage}
                    rowsPerPageOptions={[5, 10, 25]}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                />
            </Card>
            <ErrorToast
                open={errorToast.open}
                message={errorToast.message}
                onClose={handleCloseErrorToast}
            />
        </Container>
    );
}
