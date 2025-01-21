import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  Grid,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import { Category } from '../types';
import { useAppContext } from '../context/AppContext';
import { categoryApi } from '../services/api';
import { CategoryForm } from '../components/forms/CategoryForm';

export const Categories: React.FC = () => {
  const { state, dispatch } = useAppContext();
  const [openForm, setOpenForm] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await categoryApi.getAll();
        dispatch({ type: 'SET_CATEGORIES', payload: response });
      } catch (error) {
        console.error('Error fetching categories:', error);
      }
    };

    fetchCategories();
  }, [dispatch]);

  const handleAddCategory = async (data: Omit<Category, 'id' | 'createdAt' | 'updatedAt'>) => {
    try {
      const response = await categoryApi.create(data);
      dispatch({ type: 'ADD_CATEGORY', payload: response });
      setOpenForm(false);
    } catch (error) {
      console.error('Error adding category:', error);
    }
  };

  const handleUpdateCategory = async (data: Omit<Category, 'id' | 'createdAt' | 'updatedAt'>) => {
    if (!selectedCategory) return;
    
    try {
      const response = await categoryApi.update(selectedCategory.id, data);
      dispatch({ type: 'UPDATE_CATEGORY', payload: response });
      setSelectedCategory(null);
      setOpenForm(false);
    } catch (error) {
      console.error('Error updating category:', error);
    }
  };

  const handleDeleteCategory = async (id: string) => {
    try {
      await categoryApi.delete(id);
      dispatch({ type: 'DELETE_CATEGORY', payload: id });
    } catch (error) {
      console.error('Error deleting category:', error);
    }
  };

  return (
    <Box>
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4">Categories</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => {
            setSelectedCategory(null);
            setOpenForm(true);
          }}
        >
          Add Category
        </Button>
      </Box>

      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Expense Categories
              </Typography>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Subcategories</TableCell>
                      <TableCell align="right">Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {state.categories
                      .filter((category) => category.type === 'expense')
                      .map((category) => (
                        <TableRow key={category.id}>
                          <TableCell>{category.name}</TableCell>
                          <TableCell>
                            {category.subcategories?.map((sub) => (
                              <Chip
                                key={sub}
                                label={sub}
                                size="small"
                                variant="outlined"
                                sx={{ mr: 1, mb: 1 }}
                              />
                            ))}
                          </TableCell>
                          <TableCell align="right">
                            <IconButton
                              size="small"
                              onClick={() => {
                                setSelectedCategory(category);
                                setOpenForm(true);
                              }}
                            >
                              <EditIcon />
                            </IconButton>
                            <IconButton
                              size="small"
                              onClick={() => handleDeleteCategory(category.id)}
                            >
                              <DeleteIcon />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Income Categories
              </Typography>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Name</TableCell>
                      <TableCell>Subcategories</TableCell>
                      <TableCell align="right">Actions</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {state.categories
                      .filter((category) => category.type === 'income')
                      .map((category) => (
                        <TableRow key={category.id}>
                          <TableCell>{category.name}</TableCell>
                          <TableCell>
                            {category.subcategories?.map((sub) => (
                              <Chip
                                key={sub}
                                label={sub}
                                size="small"
                                variant="outlined"
                                sx={{ mr: 1, mb: 1 }}
                              />
                            ))}
                          </TableCell>
                          <TableCell align="right">
                            <IconButton
                              size="small"
                              onClick={() => {
                                setSelectedCategory(category);
                                setOpenForm(true);
                              }}
                            >
                              <EditIcon />
                            </IconButton>
                            <IconButton
                              size="small"
                              onClick={() => handleDeleteCategory(category.id)}
                            >
                              <DeleteIcon />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <CategoryForm
        open={openForm}
        onClose={() => {
          setOpenForm(false);
          setSelectedCategory(null);
        }}
        category={selectedCategory}
        onSubmit={selectedCategory ? handleUpdateCategory : handleAddCategory}
      />
    </Box>
  );
};
