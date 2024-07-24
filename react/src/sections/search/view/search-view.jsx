import { useEffect, useState } from 'react';

import {
    Stack,
    Table,
    Container,
    TableBody,
    TableRow,
    TableCell,
    TableContainer,
    TablePagination,
    Typography,
    Card,
} from '@mui/material';

import { generateMockRepos } from 'src/_mock/repo';

import Scrollbar from 'src/components/scrollbar';

import SearchForm from 'src/components/search/searchForm';
import ErrorToast from 'src/components/toast/errorToast';
import SuccessToast from 'src/components/toast/successToast';

import TableNoData from '../table-no-data';
import RepoTableRow from '../search-table-row';
import RepoTableHead from '../search-table-head';
import TableEmptyRows from '../table-empty-rows';
import RepoTableToolbar from '../search-table-toolbar';
import { emptyRows, applyFilter, getComparator } from '../utils';

import { searchGithubRepos, saveRepos } from 'src/api/client';

// ----------------------------------------------------------------------

export default function SearchPage() {
    const [page, setPage] = useState(0);
    const [order, setOrder] = useState('asc');
    const [selected, setSelected] = useState([]);
    const [orderBy, setOrderBy] = useState('name');
    const [filterName, setFilterName] = useState('');
    const [rowsPerPage, setRowsPerPage] = useState(10);

    const [isSearching, setIsSearching] = useState(false);
    const [errorToast, setErrorToast] = useState({ open: false, message: '' });
    const [successToast, setSuccessToast] = useState({ open: false, message: '' });
    const [repos, setRepos] = useState([]);
    const [searchParams, setSearchParams] = useState({
        name: '',
        query: '',
        opts: {
            sort: '',
            order: 'descending',
        },
        startPage: 0,
        currentPage: 0,
        pagesToProcess: 1,
    });

    useEffect(() => {
        // Load data from localStorage when component mounts
        const storedRepos = localStorage.getItem('storedRepos');
        if (storedRepos) {
            setRepos(JSON.parse(storedRepos));
        }
    }, [setRepos]);

    const saveToLocalStorage = (repos) => {
        localStorage.setItem('storedRepos', JSON.stringify(repos));
    };

    const handleSort = (event, id) => {
        const isAsc = orderBy === id && order === 'asc';
        if (id !== '') {
            setOrder(isAsc ? 'desc' : 'asc');
            setOrderBy(id);
        }
    };

    const handleSelectAllClick = (event) => {
        if (event.target.checked) {
            const newSelecteds = repos.map((n) => n);
            setSelected(newSelecteds);
            return;
        }
        setSelected([]);
    };

    const handleClick = (event, row, index) => {
        const selectedIndex = selected.findIndex((item) => item === row);
        let newSelected = [];

        if (selectedIndex === -1) {
            // If the row is not in the selected array, add it
            newSelected = [...selected, row];
        } else if (selectedIndex === 0) {
            // If it's the first item, remove it
            newSelected = selected.slice(1);
        } else if (selectedIndex === selected.length - 1) {
            // If it's the last item, remove it
            newSelected = selected.slice(0, -1);
        } else if (selectedIndex > 0) {
            // If it's in the middle, remove it
            newSelected = [
                ...selected.slice(0, selectedIndex),
                ...selected.slice(selectedIndex + 1),
            ];
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
        inputData: repos,
        comparator: getComparator(order, orderBy),
        filterName,
    });

    const validateForm = () => {
        // TODO: Generate UUID for name (so we can track searches)
        // OR just drop it.. undecided
        // if (!searchParams.name.trim()) {
        //     return 'Search Name is required';
        // }

        if (!searchParams.query.trim()) {
            return 'Query is required';
        }

        if (searchParams.startPage < 0) {
            return 'Start Page must be at least 0';
        }

        if (searchParams.pagesToProcess < 1) {
            return 'Pages to Process must be at least 1';
        }

        if (searchParams.pagesToProcess > 10) {
            return 'Pages to Process is limited to 20';
        }

        return null;
    };

    const handleCloseErrorToast = () => {
        setErrorToast({ open: false, message: '' });
    };

    const handleCloseSuccessToast = () => {
        setSuccessToast({ open: false, message: '' });
    };

    const handleSearch = async () => {
        const error = validateForm();

        if (error) {
            setErrorToast({ open: true, message: error });
            return;
        }

        setIsSearching(true);
        console.log('Making github api search');
        console.log(searchParams);
        console.log('Starting search...');

        try {
            const res = await searchGithubRepos(searchParams);
            if (res === null) {
                setRepos([]);
                saveToLocalStorage([]);
            } else {
                setRepos(res);
                saveToLocalStorage(res);
            }

            console.log('Search completed!');
        } catch (e) {
            setErrorToast({ open: true, message: e.response.data.message });
        }

        setSelected([]);
        setIsSearching(false);
    };

    const handleSaveRepos = async () => {
        console.log('Saving repos to db');
        console.log('Rows to save:', selected);

        try {
            await saveRepos(selected);
            setSuccessToast({ open: true, message: `${selected.length} repo(s) saved to DB.` });
        } catch (e) {
            console.log(e);
            setErrorToast({ open: true, message: e.response.data.message });
        }

        setSelected([]);
    };

    const notFound = !dataFiltered.length && !!filterName;

    return (
        <Container>
            <Stack direction="column" alignItems="left" justifyContent="space-between" mb={5}>
                <Stack mb={1}>
                    <Typography variant="h4">Search for new Repositories</Typography>
                </Stack>
                <Typography variant="p">
                    Here you can search GitHub for new repositories to add to the GooGL-Bye queue.
                </Typography>
            </Stack>

            <SearchForm
                setSearchParams={setSearchParams}
                handleSearch={handleSearch}
                searchParams={searchParams}
                isSearching={isSearching}
            />

            {}
            <Card>
                <RepoTableToolbar
                    numSelected={selected.length}
                    filterName={filterName}
                    onFilterName={handleFilterByName}
                    saveSelectedRepos={handleSaveRepos}
                />

                <Scrollbar>
                    <TableContainer sx={{ overflow: 'unset' }}>
                        <Table sx={{ minWidth: 800 }}>
                            <RepoTableHead
                                order={order}
                                orderBy={orderBy}
                                rowCount={repos.length}
                                numSelected={selected.length}
                                onRequestSort={handleSort}
                                onSelectAllClick={handleSelectAllClick}
                                headLabel={[
                                    { id: 'name', label: 'Repository Name' },
                                    { id: 'author', label: 'Author' },
                                    { id: 'language', label: 'Language' },
                                    { id: 'stars', label: 'Stars', align: 'center' },
                                    { id: 'forks', label: 'Forks', align: 'center' },
                                    { id: '' },
                                ]}
                            />
                            <TableBody>
                                {repos.length > 0 ? (
                                    <>
                                        {dataFiltered
                                            .slice(
                                                page * rowsPerPage,
                                                page * rowsPerPage + rowsPerPage
                                            )
                                            .map((row, index) => (
                                                <RepoTableRow
                                                    key={index}
                                                    name={row.name}
                                                    author={row.author}
                                                    language={row.language}
                                                    ghUrl={row.ghUrl}
                                                    stars={row.stars}
                                                    forks={row.forks}
                                                    avatarUrl={row.avatarUrl}
                                                    lastCommit={row.lastCommit}
                                                    selected={selected.some(
                                                        (selectedRow) => selectedRow === row
                                                    )}
                                                    handleClick={(event) =>
                                                        handleClick(event, row, index)
                                                    }
                                                />
                                            ))}
                                        <TableEmptyRows
                                            height={77}
                                            emptyRows={emptyRows(page, rowsPerPage, repos.length)}
                                        />
                                        {notFound && <TableNoData query={filterName} />}
                                    </>
                                ) : (
                                    <TableRow>
                                        <TableCell colSpan={7} align="center">
                                            <Typography variant="subtitle1">
                                                No repositories found.
                                            </Typography>
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Scrollbar>
                <TablePagination
                    page={page}
                    component="div"
                    count={repos.length}
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
            <SuccessToast
                open={successToast.open}
                message={successToast.message}
                onClose={handleCloseSuccessToast}
            />
        </Container>
    );
}
